package configs

import (
	"encoder/configs/video"
	"fmt"
)

type EncoderConfig struct {
	ProgressUpdate float32
	Thread         int
	Format         Formats
}

type Formats struct {
	Name      string
	Extension string
	Params    []string
}

func GetConfigWithExtension(ext, resolution string) *Formats {
	ex := fmt.Sprintf("%v_%v", ext, resolution)
	switch ext {
	case "mp4":
		mp4Configs := video.GetMp4()
		conf, ok := mp4Configs[ex]
		if ok {
			return conf
		}
	case "avi":
		aviConfigs := video.GetAvi()
		conf, ok := aviConfigs[ex]
		if ok {
			return conf
		}
	default:
		return &Formats{}
	}
	return &Formats{}
}
