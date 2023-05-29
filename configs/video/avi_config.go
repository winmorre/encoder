package video

import "encoder/configs"

var aviUtraHighDef = &configs.Formats{
	Name:      "avi_uhd",
	Extension: "avi",
	Params:    []string{},
}

var aviHighDef = &configs.Formats{
	Name:      "avi_hd",
	Extension: "avi",
	Params:    []string{},
}

var aviStandardDef = &configs.Formats{
	Name:      "avi_sd",
	Extension: "avi",
	Params:    []string{},
}

var aviLowDef = &configs.Formats{
	Name:      "avi_ld",
	Extension: "avi",
	Params:    []string{},
}
var aviUtraLowDef = &configs.Formats{
	Name:      "avi_uld",
	Extension: "avi",
	Params:    []string{},
}

func GetAvi() map[string]*configs.Formats {
	return map[string]*configs.Formats{
		"avi_uhd": aviUtraHighDef,
		"avi_hd":  aviHighDef,
		"avi_sd":  aviStandardDef,
		"avi_ld":  aviLowDef,
		"avi_uld": aviUtraLowDef,
	}
}
