package config

import (
	"strings"
)

type Config struct {
	MaxWidth     int
	MaxHeight    int
	OutputFormat string
	Quality      int // for JPEG/WebP (1-100)
}

func NewConfig(width, height int, format string, quality int) *Config {
	normFormat := strings.ToLower(format)
	if normFormat == "jpeg" {
		normFormat = "jpg"
	}
	if quality < 1 {
		quality = 1
	} else if quality > 100 {
		quality = 100
	}
	return &Config{
		MaxWidth:     width,
		MaxHeight:    height,
		OutputFormat: normFormat,
		Quality:      quality,
	}
}

func DefaultConfig() *Config {
	return NewConfig(1200, 1200, "png", 75)
}

func WebpConfig() *Config {
	return NewConfig(1200, 1200, "webp", 85)
}

func JpgConfig() *Config {
	return NewConfig(1200, 1200, "jpg", 80)
}

func (c *Config) WithMaxSize(width, height int) *Config {
	c.MaxWidth = width
	c.MaxHeight = height
	return c
}
