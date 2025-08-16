package fileio

import (
	"fmt"
	"image"
	"os"

	"github.com/disintegration/imaging"
	"github.com/kolesa-team/go-webp/encoder"
)

func HandleImageResize(img image.Image) image.Image {
	return imaging.Fit(img, 1200, 1200, imaging.Lanczos)
}

func SaveImageWebP(img image.Image, outputPath string, quality int) error {
	webpFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating WebP file: %v", err)
	}
	defer webpFile.Close()
	qualityFloat := float32(quality) / 100.0
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetPhoto, qualityFloat)
	if err != nil {
		return fmt.Errorf("error creating encoder options: %v", err)
	}
	enc, err := encoder.NewEncoder(img, options)
	if err != nil {
		return fmt.Errorf("error creating encoder: %v", err)
	}
	if err = enc.Encode(webpFile); err != nil {
		return fmt.Errorf("error encoding WebP: %v", err)
	}
	return nil
}

func SaveImagePNG(img image.Image, outputPath string) error {
	return imaging.Save(img, outputPath)
}
