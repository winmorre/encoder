package services

import (
	"context"
	"encoder/models"
	"fmt"
	"path/filepath"
	"strings"
)

func ConvertVideo(ctx context.Context, tempFileInfo chan *models.TempFileOutput, params []string, resolution, ext string, targetPath chan string, execError chan error) {
	/* Converts video file into a defined format*/
	var backend = FFmpegBackend{ffmpegPath: "ffmpeg", ffprobePath: "ffprobe"}
	tempFile := <-tempFileInfo

	execError <- tempFile.Err

	filename := strings.TrimSuffix(tempFile.Filename, filepath.Ext(tempFile.Filename))

	targetPath <- ""
	outputPath := fmt.Sprintf("temp_vids/%v_%v.%v", filename, resolution, ext)
	encode, err := backend.Encode(ctx, tempFile.Filename, outputPath, params)

	execError <- err
	fmt.Printf("File encoded at path: %v,  %v", encode, targetPath)

	for {
		select {
		case <-ctx.Done():
			targetPath <- outputPath
			execError <- ctx.Err()
			break
		default:
			targetPath <- ""
			execError <- ctx.Err()
		}
	}
}

//func createTempFileForVideo(file *os.File, path chan<- string) {
//	tempFile, err := os.CreateTemp("", file.Name())
//
//	if err != nil {
//		panic("Temp file not created")
//	}
//	vuf := make([]byte, 512)
//	_, _ = file.Read(vuf)
//	_, _ = tempFile.Write(vuf)
//	path <- tempFile.Name()
//}

/*
- see how to add gin rest
- get the uploaded file
- save temporary
- encode
- return the encoded file
- remove temp file
*/
