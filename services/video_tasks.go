package services

import (
	"encoder/configs"
	"fmt"
	"os"
)

func ConvertVideo(file os.File, format configs.Formats) {
	/* Converts video file into a defined format*/
	fileName := make(chan string)

	go createTempFileForVideo(&file, fileName)
	var backend FFmpegBackend

	encode, err := backend.Encode(<-fileName, "", format.Params)

	if err != nil {
		panic("Error encoding")
	}
	fmt.Printf("File encoded: %v", <-encode)
}

//func encode(srcPath string, format configs.Formats, options map[string]string) {
//	// create a temp file, if not already existed
//	// encode this
//}

func createTempFileForVideo(file *os.File, path chan<- string) {
	tempFile, err := os.CreateTemp("", file.Name())

	if err != nil {
		panic("Temp file not created")
	}
	vuf := make([]byte, 512)
	_, _ = file.Read(vuf)
	_, _ = tempFile.Write(vuf)
	path <- tempFile.Name()
}