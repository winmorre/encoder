package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const regexTimeCode = `"time=(\d+:\d+.\d+)"`

type FFmpegBackend struct {
	Params      []string
	ffmpegPath  string
	ffprobePath string
}

func spawn(cmd []string) *exec.Cmd {
	ffmpegCmd := exec.Command(cmd[0], cmd[1:]...)
	return ffmpegCmd
}

func parseStreamOutput(data string) (*StreamOutput, error) {
	mediaInfo := StreamOutput{}

	err := json.Unmarshal([]byte(data), &mediaInfo)

	if err == nil {
		var videoStream = make([]Stream, 0, len(mediaInfo.Streams))
		var audioStream = make([]Stream, 0, len(mediaInfo.Streams))
		var subtitleStream = make([]Stream, 0, len(mediaInfo.Streams))
		for _, o := range mediaInfo.Streams {
			if o.CodecType == "video" {
				videoStream = append(videoStream, o)
			} else if o.CodecType == "audio" {
				audioStream = append(audioStream, o)
			} else if o.CodecType == "subtitle" {
				subtitleStream = append(subtitleStream, o)
			}
		}

		mediaInfo.Audio = audioStream
		mediaInfo.Subtitle = subtitleStream
		mediaInfo.Video = videoStream
	}

	mediaInfo.Streams = nil
	return &mediaInfo, err
}

func (ff *FFmpegBackend) GetMediaInfo(ctx context.Context, videoPath string, outputInfo chan *StreamOutput, execError chan error) {
	cmd := []string{ff.ffprobePath, "-i", videoPath, "-hide_banner", "-loglevel", "warning", "-of", "json", "-show_format", "-show_streams"}
	execCmd := spawn(cmd)
	var mediaInfo *StreamOutput

	var builder strings.Builder
	execCmd.Stdout = &builder
	err := execCmd.Run()
	execError <- err
	//if err != nil {
	//	fmt.Printf("Error occurred while getting media Info:  %v", err.Error())
	//	return
	//}
	mediaInfo, err = parseStreamOutput(builder.String())

	//if err != nil {
	//	fmt.Printf("Error occurred parsing streamOuput:  %v", err.Error())
	//	return
	//}
	outputInfo <- mediaInfo
	execError <- err
}

func (ff *FFmpegBackend) GetThumbnail(ctx context.Context, videoPath string, atTime float64) string {
	// Extract and image from a video
	fileName := path.Base(videoPath)
	fileNameSplit := strings.Split(fileName, ".")
	fileName = strings.Join(fileNameSplit[:len(fileNameSplit)-1], " ")
	file, err := os.CreateTemp("", fmt.Sprintf("%v.jpg", fileName))
	if err != nil {
		panic(fmt.Sprintf("Temp file not creted %v", fileName))
	}
	mediaInfo := make(chan *StreamOutput)
	execError := make(chan error)

	defer close(mediaInfo)
	defer close(execError)

	go ff.GetMediaInfo(ctx, videoPath, mediaInfo, execError)

	for {
		select {
		case <-ctx.Done():
			break
		case err := <-execError:
			if err != nil {
				break
			}
		case media := <-mediaInfo:
			if len(media.Video) != 0 {
				return getMediaFilePath(ff, file, media, fileName, videoPath, atTime)
			}
		}
	}
}

func getMediaFilePath(backend *FFmpegBackend, file *os.File, mediaInfo *StreamOutput, fileName, videoPath string, atTime float64) string {
	duration, _ := strconv.ParseFloat(mediaInfo.Format.Duration, 64)
	var err error
	if atTime > duration {
		panic(fmt.Sprintf("atTime (%v) is greater than the video duration", atTime))
	}
	defer file.Close()

	var fPath string
	fPath, err = filepath.Abs(file.Name())
	if err == nil {
		panic(fmt.Sprintf("Error occured getting Abs path: %v", err.Error()))
	}

	cmd := []string{backend.ffmpegPath, "-i", videoPath, "-vframes", "1", "-ss", fmt.Sprintf("%v", atTime), "-y", fPath}

	execCmd := spawn(cmd)
	err = execCmd.Run()

	if err != nil {
		panic(fmt.Sprintf("ffmpeg could not create thumnail:  %v", err.Error()))
	}

	var imageFileInfo os.FileInfo
	imageFileInfo, err = file.Stat()

	if err != nil {
		panic(fmt.Sprintf("Could not get file stat:  %v", err.Error()))
	}

	if imageFileInfo.Size() == 0.0 {
		os.Remove(fileName)
		panic("Failed to generate thumbnail")
	}
	return fPath
}

func (ff *FFmpegBackend) Encode(ctx context.Context, srcPath, targetPath string, params []string) (float64, error) {
	/* Encode a video*/
	mediaInfo := make(chan *StreamOutput)
	execError := make(chan error)
	defer close(execError)
	defer close(mediaInfo)

	go ff.GetMediaInfo(ctx, srcPath, mediaInfo, execError)
	var totalDuration string

	select {
	case <-ctx.Done():
	case media := <-mediaInfo:
		totalDuration = media.Format.Duration
	}

	totalPercentage := make(chan float64)

	defer close(totalPercentage)
	cmd := []string{ff.ffmpegPath, "-i", srcPath}
	cmd = append(cmd, params...)
	cmd = append(cmd, targetPath)

	execCmd := spawn(cmd)
	var errBuilder strings.Builder
	execCmd.Stderr = &errBuilder
	var outPutBuilder strings.Builder
	execCmd.Stdout = &outPutBuilder

	go func() {
		exErr := execCmd.Run()
		execError <- exErr
		k, _ := regexp.Compile(regexTimeCode)
		timeStr := k.FindAllString(errBuilder.String(), -1)[0]

		var Time float64
		for _, t := range strings.Split(timeStr, ":") {
			parseT, err := strconv.ParseFloat(t, 64)
			if err != nil {
				panic(fmt.Sprintf("Error converting timeStr;  %v", err.Error()))
			}
			Time = 60 * parseT
		}
		totalDurationF, err := strconv.ParseFloat(totalDuration, 64)
		if err != nil {
			panic(fmt.Sprintf("Error converting totalDuration;  %v", err.Error()))
		}
		percent := math.Round(Time / totalDurationF)
		totalPercentage <- percent

		fmt.Printf("Percentage return %v", percent)

		if m, _ := os.Stat(targetPath); m.Size() == 0 {
			panic("File size of generated file is 0")
		}

		if execCmd.ProcessState.ExitCode() != 0 {
			panic(fmt.Sprintf("%v exited with code %v", execCmd.Args, execCmd.ProcessState.ExitCode()))
		}
	}()

	for {
		select {
		case <-ctx.Done():
			totalPercentage <- 0.0
			execError <- ctx.Err()
			return 0.0, ctx.Err()
		default:
			return <-totalPercentage, <-execError

		}
	}
}
