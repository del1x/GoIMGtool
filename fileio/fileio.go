package fileio

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/kolesa-team/go-webp/encoder"
)

func ReadDir(dir string) ([]os.DirEntry, error) {
	return os.ReadDir(dir)
}

func CreateDir(dir string) error {
	err := os.Mkdir(dir, 0755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("error creating directory: %v", err)
	}
	return nil
}

func LoadImage(path string) (image.Image, error) {
	img, err := imaging.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error loading image: %v", err)
	}
	return img, nil
}

func SaveImage(img image.Image, outputPath, outputFormat string) error {
	base := strings.TrimSuffix(outputPath, filepath.Ext(outputPath))
	if outputFormat == "webp" {
		outputPath = base + ".webp"
		webpFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("error creating WebP file: %v", err)
		}
		defer webpFile.Close()
		options, err := encoder.NewLosslessEncoderOptions(encoder.PresetDefault, 6)
		if err != nil {
			return fmt.Errorf("error creating encoder options: %v", err)
		}
		enc, err := encoder.NewEncoder(img, options)
		if err != nil {
			return fmt.Errorf("error creating encoder: %v", err)
		}
		err = enc.Encode(webpFile)
		if err != nil {
			return fmt.Errorf("error encoding WebP: %v", err)
		}
	} else {
		outputPath = base + ".png"
		err := imaging.Save(img, outputPath)
		if err != nil {
			return fmt.Errorf("error saving image: %v", err)
		}
	}
	fmt.Printf("Processed image: %s\n", outputPath)
	return nil
}
