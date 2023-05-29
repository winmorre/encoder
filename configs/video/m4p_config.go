package video

import "encoder/configs"

var m4pUltraHighDef = &configs.Formats{
	Name:      "m4p_uhd",
	Extension: "m4p",
	Params:    []string{},
}

var m4pHighDef = &configs.Formats{
	Name:      "m4p_hd",
	Extension: "m4p",
	Params:    []string{},
}

var m4pStandardDef = &configs.Formats{
	Name:      "m4p_sd",
	Extension: "m4p",
	Params:    []string{},
}

var m4pLowDef = &configs.Formats{
	Name:      "m4p_ld",
	Extension: "m4p",
	Params:    []string{},
}

var m4pUltraLowDef = &configs.Formats{
	Name:      "m4p_uld",
	Extension: "m4p",
	Params:    []string{},
}

func GetMp4() map[string]*configs.Formats {
	return map[string]*configs.Formats{
		"m4p_uhd": m4pUltraHighDef,
		"m4p_hd":  m4pHighDef,
		"m4p_sd":  m4pStandardDef,
		"m4p_ld":  m4pLowDef,
		"m4p_uld": m4pUltraLowDef,
	}
}
