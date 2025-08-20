package fileio

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"image/draw"

	"github.com/del1x/GoIMGtool/config"
	"github.com/disintegration/imaging"
	"github.com/kolesa-team/go-webp/encoder"
)

func HandleImageResize(img image.Image, cfg *config.Config) image.Image {
	return imaging.Fit(img, cfg.MaxWidth, cfg.MaxHeight, imaging.Lanczos)
}

func SaveImageWebP(img image.Image, outputPath string, quality int) error {
	webpFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating WebP file: %v", err)
	}
	defer webpFile.Close()
	qualityFloat := float32(quality) // Убрано деление на 100
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
	fmt.Println("Starting PNG encoding for:", outputPath)
	pngFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating PNG file: %v", err)
	}
	defer pngFile.Close()

	fmt.Printf("Input image type: %T\n", img)

	imgNRGBA := image.NewNRGBA(img.Bounds())
	draw.Draw(imgNRGBA, imgNRGBA.Bounds(), img, image.Point{0, 0}, draw.Over)
	fmt.Println("Image converted to NRGBA")

	if err := png.Encode(pngFile, imgNRGBA); err != nil {
		return fmt.Errorf("error encoding PNG %s: %v", outputPath, err)
	}
	fmt.Println("PNG encoding successful for:", outputPath)
	return nil
}
