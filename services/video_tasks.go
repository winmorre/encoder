package services

import (
	"encoder/configs"
	"encoder/models"
	"fmt"
	"os"
)

func ConvertVideo(tempFileInfo <-chan models.TempFileOutput, format configs.Formats) {
	/* Converts video file into a defined format*/
	fileName := make(chan string)

	var backend = FFmpegBackend{ffmpegPath: "ffmpeg", ffprobePath: "ffprobe"}

	encode, err := backend.Encode(<-fileName, "", format.Params)

	if err != nil {
		panic("Error encoding")
	}
	fmt.Printf("File encoded: %v", <-encode)
}

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

/*
- see how to add gin rest
- get the uploaded file
- save temporary
- encode
- return the encoded file
- remove temp file
*/
