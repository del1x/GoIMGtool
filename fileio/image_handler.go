package fileio

import (
	"image"

	"github.com/del1x/GoIMGtool/config"
	"github.com/disintegration/imaging"
)

func HandleImageResize(img image.Image, cfg *config.Config) image.Image {
	return imaging.Fit(img, cfg.MaxWidth, cfg.MaxHeight, imaging.Lanczos)
}
