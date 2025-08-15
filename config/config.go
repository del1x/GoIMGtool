package config

type Config struct {
	MaxWidth     int
	MaxHeight    int
	OutputFormat string
	Quality      int // Для JPEG/WebP
}

func DefaultConfig() *Config {
	return &Config{
		MaxWidth:     1200,
		MaxHeight:    1200,
		OutputFormat: "png",
		Quality:      75,
	}
}
