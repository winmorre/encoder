package services

type BaseEncodingBackend interface {
	Encode(srcPath, targetPath string, params []string) interface{}

	GetMediaInfo(videoPath string) *StreamOutput
	GetThumbnail(videoPath string, atTime float64) string
}

type Stream struct {
	Index          int    `json:"index"`
	CodecName      string `json:"codec_name"`
	CodecLongName  string `json:"codec_long_name"`
	Profile        string `json:"profile"`
	CodecType      string `json:"codec_type"`
	CodecTagString string `json:"codec_tag_string"`
	CodecTag       string `json:"codec_tag"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	CodedWidth     int    `json:"coded_width"`
	CodedHeight    int    `json:"coded_height"`
	ClosedCaptions int    `json:"closed_captions"`
	FilmGrain      int    `json:"film_grain"`
	HasBFrames     int    `json:"has_b_frames"`
	PixFmt         string `json:"pix_fmt"`
	Level          int    `json:"level"`
	ColorRange     string `json:"color_range"`
	ColorSpace     string `json:"color_space"`
	ColorTransfer  string `json:"color_transfer"`
	ColorPrimaries string `json:"color_primaries"`
	Refs           int    `json:"refs"`
	Id             string `json:"id"`
	RFrameRate     string `json:"r_frame_rate"`
	AvgFrameRate   string `json:"avg_frame_rate"`
	TimeBase       string `json:"time_base"`
	StartPts       int    `json:"start_pts"`
	StartTime      string `json:"start_time"`
	DurationTs     int    `json:"duration_ts"`
	Duration       string `json:"duration"`
	BitRate        string `json:"bit_rate"`
	NbFrames       string `json:"nb_frames"`
	ExtradataSize  int    `json:"extradata_size"`
	Disposition    struct {
		Default         int `json:"default"`
		Dub             int `json:"dub"`
		Original        int `json:"original"`
		Comment         int `json:"comment"`
		Lyrics          int `json:"lyrics"`
		Karaoke         int `json:"karaoke"`
		Forced          int `json:"forced"`
		HearingImpaired int `json:"hearing_impaired"`
		VisualImpaired  int `json:"visual_impaired"`
		CleanEffects    int `json:"clean_effects"`
		AttachedPic     int `json:"attached_pic"`
		TimedThumbnails int `json:"timed_thumbnails"`
		Captions        int `json:"captions"`
		Descriptions    int `json:"descriptions"`
		Metadata        int `json:"metadata"`
		Dependent       int `json:"dependent"`
		StillImage      int `json:"still_image"`
	} `json:"disposition"`
	Tags struct {
		Language    string `json:"language"`
		HandlerName string `json:"handler_name"`
		VendorId    string `json:"vendor_id"`
	} `json:"tags"`
}

type StreamFormat struct {
	Filename       string `json:"filename"`
	NbStreams      int    `json:"nb_streams"`
	NbPrograms     int    `json:"nb_programs"`
	FormatName     string `json:"format_name"`
	FormatLongName string `json:"format_long_name"`
	StartTime      string `json:"start_time"`
	Duration       string `json:"duration"`
	Size           string `json:"size"`
	BitRate        string `json:"bit_rate"`
	ProbeScore     int    `json:"probe_score"`
	Tags           struct {
		MajorBrand       string `json:"major_brand"`
		MinorVersion     string `json:"minor_version"`
		CompatibleBrands string `json:"compatible_brands"`
		Title            string `json:"title"`
		Artist           string `json:"artist"`
		Date             string `json:"date"`
		Encoder          string `json:"encoder"`
		Comment          string `json:"comment"`
	} `json:"tags"`
}

type StreamOutput struct {
	Streams  []Stream     `json:"streams"`
	Format   StreamFormat `json:"format"`
	Video    []Stream
	Audio    []Stream
	Subtitle []Stream
}