package controllers

import (
	"context"
	"encoder/configs"
	"encoder/models"
	"encoder/services"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

const VideoSizeLimit int64 = 100000000

var DefaultErrors = map[string]error{
	"fileSizeExceedLimit": errors.New(fmt.Sprintf("Currently supported file size for videos is %v", VideoSizeLimit)),
	"encodingError":       errors.New("encoding not Successful"),
}

func EncodeVideo(ctx *gin.Context) {
	duration := 2 * time.Minute
	newCtx, cancelFunc := context.WithTimeout(ctx, duration)

	defer cancelFunc()

	var videoReq models.Video
	err := ctx.ShouldBind(&videoReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": "File not uploaded"})
		return
	}

	if !checkFileSize(videoReq.File.Size) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": DefaultErrors["fileSizeExceedLimit"]})
		return
	}

	var tempFileOutput = make(chan *models.TempFileOutput)
	go saveFileAsTempFile(videoReq.File, tempFileOutput)

	params := configs.GetConfigWithExtension(videoReq.Extension, videoReq.Resolution)

	targetPath := make(chan string)
	execError := make(chan error)
	defer close(targetPath)
	defer close(execError)

	go services.ConvertVideo(newCtx, tempFileOutput, params.Params, videoReq.Resolution, videoReq.Extension, targetPath, execError)

	for {
		select {
		case <-newCtx.Done():
			err := <-execError
			if err != nil {
				ctx.JSON(http.StatusExpectationFailed, gin.H{"status": "failed", "msg": DefaultErrors["encodingError"]})
				return
			}
		case err := <-execError:
			if err != nil {
				ctx.JSON(http.StatusExpectationFailed, gin.H{"status": "failed", "msg": DefaultErrors["encodingError"]})
				return
			}
		case path := <-targetPath:
			if path != "" {
				// get the save generated file and return that
				// remove temp files
				var fileStat os.FileInfo
				file, err := getEncodedFile(path)
				fileStat, err = file.Stat()

				if err != nil {
					ctx.JSON(http.StatusExpectationFailed, gin.H{})
					return
				}
				fileHeader := multipart.FileHeader{
					Filename: fileStat.Name(),
					Size:     fileStat.Size(),
				}
				ctx.FileAttachment(path, fileStat.Name())
				ctx.JSON(http.StatusOK, gin.H{"status": "success", "file_header": fileHeader})
				os.Remove(path)
				tempFilePath := <-tempFileOutput
				if tempFilePath.Filename != "" {
					os.Remove(tempFilePath.Filename)
				}
				return
			}
		}
	}
}

func checkFileSize(size int64) bool {
	return size <= VideoSizeLimit
}

func saveFileAsTempFile(file *multipart.FileHeader, op chan<- *models.TempFileOutput) {
	f, err := file.Open()
	var output *models.TempFileOutput
	var fileData []byte
	var tempFile *os.File
	defer f.Close()

	if err != nil {
		output = &models.TempFileOutput{Success: false, Filename: "", Err: err}
		op <- output
		return
	}

	fileData, err = io.ReadAll(f)
	if err != nil {
		output = &models.TempFileOutput{Success: false, Filename: "", Err: err}
		op <- output
		return
	}

	tempFile, err = os.CreateTemp("temp-vids", file.Filename)
	defer tempFile.Close()

	_, err = tempFile.Write(fileData)

	if err != nil {
		output = &models.TempFileOutput{Success: false, Filename: "", Err: err}
		op <- output
		return
	}

	output = &models.TempFileOutput{Success: true, Filename: fmt.Sprintf("temp_vids/%v", file.Filename)}
	op <- output
}

func getEncodedFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	return file, err
}

func GetMediaInfo(ctx *gin.Context) {
	var request models.MediaInfoRequest
	newCtx, cancelFunc := context.WithTimeout(ctx, time.Second)

	defer cancelFunc()
	var err error
	err = ctx.ShouldBind(&request)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": "Data format not correct"})
		fmt.Printf("data binding error:  %v", err.Error())
		return
	}

	mediaInfo := make(chan *services.StreamOutput)
	execError := make(chan error)
	tempFileChan := make(chan *models.TempFileOutput)

	defer close(mediaInfo)
	defer close(execError)
	defer close(tempFileChan)

	err = nil
	saveFileAsTempFile(request.Video, tempFileChan)

	var tempInfo *models.TempFileOutput

	ff := services.FFmpegBackend{FfmpegPath: "ffmpeg", FfprobePath: "ffprobe"}

	for {
		select {
		case <-newCtx.Done():
			fmt.Printf("New context Done:  %v", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err.Error()})
			break
		case tempInfo = <-tempFileChan:
			if tempInfo.Success {
				go ff.GetMediaInfo(ctx, tempInfo.Filename, mediaInfo, execError)
			}
		case err = <-execError:
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err.Error()})
				break
			}
		case media := <-mediaInfo:
			if media != nil {
				ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": media})
				goto removeTempFile
			}
		}
	}

removeTempFile:
	{
		defer removeFile(tempInfo.Filename)
		return
	}
}
func removeFile(filePath string) {
	defer os.Remove(filePath)
}
