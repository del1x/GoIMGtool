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

	img = HandleImageResize(img)
	fmt.Println("Image resized, type:", fmt.Sprintf("%T", img))

	bestQuality, err := OptimizeQuality(img, outputFormat, base, p.TargetSizeKB)
	if err != nil {
		return fmt.Errorf("error optimizing quality: %v", err)
	}
	fmt.Println("Optimized quality:", bestQuality)

	if outputFormat == "webp" {
		outputPath = base + ".webp"
		if err := SaveImageWebP(img, outputPath, bestQuality); err != nil {
			return err
		}
	} else if outputFormat == "jpg" {
		outputPath = base + ".jpg"
		if err := ManageExif(img, outputPath, bestQuality); err != nil {
			return err
		}
	} else if outputFormat == "png" {
		outputPath = base + ".png"
		if err := SaveImagePNG(img, outputPath); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupported format: %s", outputFormat)
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
