package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name    string
		width   int
		height  int
		format  string
		quality int
		want    *Config
	}{
		{
			name:    "Default PNG",
			width:   1200,
			height:  1200,
			format:  "png",
			quality: 75,
			want:    &Config{MaxWidth: 1200, MaxHeight: 1200, OutputFormat: "png", Quality: 75},
		},
		{
			name:    "JPEG to JPG",
			width:   800,
			height:  800,
			format:  "jpeg",
			quality: 90,
			want:    &Config{MaxWidth: 800, MaxHeight: 800, OutputFormat: "jpg", Quality: 90},
		},
		{
			name:    "Quality Limit Low",
			width:   500,
			height:  500,
			format:  "webp",
			quality: 0,
			want:    &Config{MaxWidth: 500, MaxHeight: 500, OutputFormat: "webp", Quality: 1},
		},
		{
			name:    "Quality Limit High",
			width:   1500,
			height:  1500,
			format:  "jpg",
			quality: 150,
			want:    &Config{MaxWidth: 1500, MaxHeight: 1500, OutputFormat: "jpg", Quality: 100},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewConfig(tt.width, tt.height, tt.format, tt.quality)
			if got.MaxWidth != tt.want.MaxWidth ||
				got.MaxHeight != tt.want.MaxHeight ||
				got.OutputFormat != tt.want.OutputFormat ||
				got.Quality != tt.want.Quality {
				t.Errorf("NewConfig() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	want := &Config{MaxWidth: 1200, MaxHeight: 1200, OutputFormat: "png", Quality: 75}
	if got := DefaultConfig(); got.MaxWidth != want.MaxWidth ||
		got.MaxHeight != want.MaxHeight ||
		got.OutputFormat != want.OutputFormat ||
		got.Quality != want.Quality {
		t.Errorf("DefaultConfig() = %+v, want %+v", got, want)
	}
}

func TestWebpConfig(t *testing.T) {
	want := &Config{MaxWidth: 1200, MaxHeight: 1200, OutputFormat: "webp", Quality: 85}
	if got := WebpConfig(); got.MaxWidth != want.MaxWidth ||
		got.MaxHeight != want.MaxHeight ||
		got.OutputFormat != want.OutputFormat ||
		got.Quality != want.Quality {
		t.Errorf("WebpConfig() = %+v, want %+v", got, want)
	}
}

func TestJpgConfig(t *testing.T) {
	want := &Config{MaxWidth: 1200, MaxHeight: 1200, OutputFormat: "jpg", Quality: 80}
	if got := JpgConfig(); got.MaxWidth != want.MaxWidth ||
		got.MaxHeight != want.MaxHeight ||
		got.OutputFormat != want.OutputFormat ||
		got.Quality != want.Quality {
		t.Errorf("JpgConfig() = %+v, want %+v", got, want)
	}
}

func TestWithMaxSize(t *testing.T) {
	cfg := DefaultConfig()
	newCfg := cfg.WithMaxSize(800, 600)
	if newCfg.MaxWidth != 800 || newCfg.MaxHeight != 600 ||
		newCfg.OutputFormat != "png" || newCfg.Quality != 75 {
		t.Errorf("WithMaxSize() = %+v, want MaxWidth=800, MaxHeight=600", newCfg)
	}
}
