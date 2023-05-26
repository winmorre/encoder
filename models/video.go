package models

import "mime/multipart"

type Video struct {
	Extension string                `form:"extension"` // format to convert to
	Name      string                `form:"Name"`
	File      *multipart.FileHeader `form:"file"`
}

func (v *Video) CreateTempFile() {

}

type TempFileOutput struct {
	Path    string
	Err     error
	Success bool
}
