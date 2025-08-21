package fileio

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/del1x/GoIMGtool/config"
)

type ImageProcessor struct {
	TargetSizeKB int
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{
		TargetSizeKB: 100,
	}
}

func (p *ImageProcessor) SaveImage(img image.Image, outputPath, outputFormat string, cfg *config.Config) error {
	fmt.Println("Processing image with format:", outputFormat)
	base := strings.TrimSuffix(outputPath, filepath.Ext(outputPath))
	outputPath = base + "." + outputFormat

	img = HandleImageResize(img, cfg)
	fmt.Println("Image resized, type:", fmt.Sprintf("%T", img))

	bestQuality, err := OptimizeQuality(img, outputFormat, base, p.TargetSizeKB)
	if err != nil {
		return fmt.Errorf("error optimizing quality: %v", err)
	}
	fmt.Println("Optimized quality:", bestQuality)

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	encoder, err := GetEncoder(outputFormat)
	if err != nil {
		return err
	}
	if err := encoder.Encode(img, file, bestQuality); err != nil {
		return fmt.Errorf("error encoding image: %v", err)
	}

	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		return fmt.Errorf("error getting file info: %v", err)
	}
	fileSizeKB := fileInfo.Size() / 1024
	fmt.Printf("Processed image: %s with size %d KB and quality: %d\n", outputPath, fileSizeKB, bestQuality)
	if int64(fileSizeKB) > int64(p.TargetSizeKB) {
		return fmt.Errorf("final size %d KB exceeds %d KB limit", fileSizeKB, p.TargetSizeKB)
	}

	return nil
}
