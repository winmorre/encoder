package configs

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
