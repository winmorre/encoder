package controllers

import (
	"encoder/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

const VideoSizeLimit float64 = 100 * 60

var DefaultErrors = map[string]error{
	"fileSizeExceedLimit": errors.New(fmt.Sprintf("Currently supported file size for videos is %v", VideoSizeLimit)),
	"encodingError":       errors.New("encoding not Successful"),
}

func EncodeVideo(ctx *gin.Context) {
	var videoReq models.Video
	err := ctx.ShouldBind(&videoReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": "File not uploaded"})
		return
	}

	if !checkFileSize(float64(videoReq.File.Size)) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": DefaultErrors["fileSizeExceedLimit"]})
		return
	}

	var tempFileOutput = make(chan *models.TempFileOutput)
	go saveFileAsTempFile(videoReq.File, tempFileOutput)

	//go services.ConvertVideo(&videoReq.File)
}

func checkFileSize(size float64) bool {
	return size <= VideoSizeLimit
}

func saveFileAsTempFile(file *multipart.FileHeader, op chan<- *models.TempFileOutput) {
	f, err := file.Open()
	var output *models.TempFileOutput
	var fileData []byte
	var tempFile *os.File
	defer f.Close()

	if err != nil {
		output = &models.TempFileOutput{Success: false, Path: "", Err: err}
		op <- output
		return
	}

	fileData, err = io.ReadAll(f)
	if err != nil {
		output = &models.TempFileOutput{Success: false, Path: "", Err: err}
		op <- output
	}

	tempFile, err = os.CreateTemp("temp-vids", file.Filename)
	defer tempFile.Close()

	_, err = tempFile.Write(fileData)

	if err != nil {
		output = &models.TempFileOutput{Success: false, Path: "", Err: err}
		op <- output
	}

	output = &models.TempFileOutput{Success: true, Path: fmt.Sprintf("temp_vids/%v", file.Filename)}
	op <- output
	close(op)
}
