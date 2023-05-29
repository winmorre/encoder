package models

import "mime/multipart"

type Video struct {
	Extension  string                `form:"extension"` // format to convert to
	Name       string                `form:"Name"`
	Resolution string                `form:"resolution"`
	File       *multipart.FileHeader `form:"file"`
}

func (v *Video) CreateTempFile() {

}

type TempFileOutput struct {
	Filename string
	Err      error
	Success  bool
}
