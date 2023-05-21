package video

import "encoder/configs"

var webmStandardDef = configs.Formats{
	Name:      "webm_sd",
	Extension: "webm",
	Params: []string{
		"-b:v", "1000k",
		"-maxrate", "1000k",
		"-bufsize", "2000k",
		"-codec:v", "libvpx",
		"-r", "30",
		"-vf", "scal=-1:480",
		"-qmin", "10",
		"-qmax", "42",
		"-codec:a", "libvorbis",
		"-b:a", "128k",
		"-f", "webm",
	},
}

var webmHighDef = configs.Formats{
	Name:      "webm_hd",
	Extension: "webm",
	Params: []string{
		"-codec:v", "libvpx",
		"-b:v", "3000k",
		"-maxrate", "3000k",
		"-bufsize", "6000k",
		"-vf", "scale=-1:720",
		"-qmin", "11",
		"-qmax", "51",
		"-acodec", "libvorbis",
		"-b:a", "128k",
		"-f", "webm",
	},
}

var webmLowDef = configs.Formats{
	Name:      "webm_ld",
	Extension: "webm",
	Params: []string{
		"-codec:v", "libvpx",
		"-b:v", "800k",
		"-maxrate", "800k",
		"-bufsize", "1000k",
		"-vf", "scale=-1:360",
		"-qmin", "11", //change this
		"-qmax", "51", // change this
		"-acodec", "libvorbis",
		"-b:a", "128k",
		"-f", "webm",
	},
}

var webmUtraLowDef = configs.Formats{
	Name:      "webm_uld",
	Extension: "webm",
	Params: []string{
		"-codec:v", "libvpx",
		"-b:v", "350k",
		"-maxrate", "350k",
		"-bufsize", "600k",
		"-vf", "scale=-1:120",
		"-qmin", "11", //change this
		"-qmax", "51", // change this
		"-acodec", "libvorbis",
		"-b:a", "128k",
		"-f", "webm",
	},
}
