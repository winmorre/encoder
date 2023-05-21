package video

import "encoder/configs"

var oggStandardDef = configs.Formats{
	Name:      "ogg_sd",
	Extension: "ogg",
	Params: []string{
		"-codec:v", "",
	},
}

var oggHighDef = configs.Formats{
	Name:      "ogg_hd",
	Extension: "ogg",
	Params: []string{
		"",
	},
}

var oggLowDef = configs.Formats{
	Name:      "ogg_hd",
	Extension: "ogg",
	Params: []string{
		"",
	},
}

var oggUltraLowDef = configs.Formats{
	Name:      "ogg_hd",
	Extension: "ogg",
	Params: []string{
		"",
	},
}
