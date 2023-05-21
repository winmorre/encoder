package video

import "encoder/configs"

var mp4StandardDef = configs.Formats{
	Name:      "mp4_sd",
	Extension: "mp4",
	Params: []string{
		"-codec:v", "libx264",
		"-crf", "20",
		"-preset", "medium",
		"-b:v", "1000k",
		"-maxrate", "1000k",
		"-bufsize", "2000k",
		"-vf", "scale=-2:480", // http://superuser.com/a/776254
		"-codec:a", "aac",
		"-b:a", "128k",
		"-strict", "-2",
	},
}

var mp4HighDef = configs.Formats{
	Name:      "mp4_hd",
	Extension: "mp4",
	Params: []string{
		"-codec:v", "libx264",
		"-crf", "20",
		"-preset", "medium",
		"-b:v", "3000k",
		"-maxrate", "3000k",
		"-bufsize", "6000k",
		"-vf", "scale=-2:720",
		"-codec:a", "aac",
		"-b:a", "128k",
		"-strict", "-2",
	},
}

var mp4LowDef = configs.Formats{
	Name:      "mp4_ld",
	Extension: "mp4",
	Params: []string{
		"-codec:v", "libx264",
		"-crf", "20",
		"-preset", "medium",
		"-b:v", "800k",
		"-maxrate", "800k",
		"-bufsize", "1000k",
		"-vf", "scale=-2:360",
		"-codec:a", "aac",
		"-b:a", "128k",
		"-strict", "-2",
	},
}

var mp4UtraLowDef = configs.Formats{
	Name:      "mp4_uld",
	Extension: "mp4",
	Params: []string{
		"-codec:v", "libx264",
		"-crf", "20",
		"-preset", "medium",
		"-b:v", "350k",
		"-maxrate", "350k",
		"-bufsize", "600k",
		"-vf", "scale=-2:120",
		"-codec:a", "aac",
		"-b:a", "128k",
		"-strict", "-2",
	},
}
