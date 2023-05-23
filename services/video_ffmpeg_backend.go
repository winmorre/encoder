package services

import (
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

func (ff *FFmpegBackend) GetMediaInfo(videoPath string) (mediaInfo *StreamOutput) {
	cmd := []string{ff.ffprobePath, "-i", videoPath, "-hide_banner", "-loglevel", "warning", "-of", "json", "-show_format", "-show_streams"}
	execCmd := spawn(cmd)

	var builder strings.Builder
	execCmd.Stdout = &builder
	err := execCmd.Run()

	if err != nil {
		fmt.Printf("Error occurred while getting media Info:  %v", err.Error())
	}
	mediaInfo, err = parseStreamOutput(builder.String())

	if err != nil {
		fmt.Printf("Error occurred parsing streamOuput:  %v", err.Error())
	}
	return mediaInfo
}

func (ff *FFmpegBackend) GetThumbnail(videoPath string, atTime float64) string {
	// Extract and image from a video
	fileName := path.Base(videoPath)
	fileNameSplit := strings.Split(fileName, ".")
	fileName = strings.Join(fileNameSplit[:len(fileNameSplit)-1], " ")
	file, err := os.CreateTemp("", fmt.Sprintf("%v.jpg", fileName))
	if err != nil {
		panic(fmt.Sprintf("Temp file not creted %v", fileName))
	}
	mediaInfo := ff.GetMediaInfo(videoPath)
	duration, _ := strconv.ParseFloat(mediaInfo.Format.Duration, 64)

	if atTime > duration {
		panic(fmt.Sprintf("atTime (%v) is greater than the video duration", atTime))
	}
	defer file.Close()

	var fPath string
	fPath, err = filepath.Abs(file.Name())
	if err == nil {
		panic(fmt.Sprintf("Error occured getting Abs path: %v", err.Error()))
	}

	cmd := []string{ff.ffmpegPath, "-i", videoPath, "-vframes", "1", "-ss", fmt.Sprintf("%v", atTime), "-y", fPath}

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

	if imageFileInfo.Size() < 1.0 {
		os.Remove(fileName)
		panic("Failed to generate thumnail")
	}
	return fPath
}

func (ff *FFmpegBackend) Encode(srcPath, targetPath string, params []string) (chan float64, chan error) {
	/* Encode a video*/
	mediaInfo := ff.GetMediaInfo(srcPath)
	totalDuration := mediaInfo.Format.Duration
	totalPercentage := make(chan float64)
	execError := make(chan error)

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

	return totalPercentage, execError
}
