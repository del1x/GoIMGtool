package fileio

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/kolesa-team/go-webp/encoder"
)

func saveAndGetSize(img image.Image, quality int, format, path string) (int, error) {
	var file *os.File
	var err error
	if format == "jpg" {
		file, err = os.Create(path)
		if err != nil {
			return 0, fmt.Errorf("error creating JPEG file: %v", err)
		}
		defer file.Close()
		if err := jpeg.Encode(file, img, &jpeg.Options{Quality: quality}); err != nil {
			return 0, fmt.Errorf("error encoding JPEG: %v", err)
		}
	} else if format == "webp" {
		file, err = os.Create(path)
		if err != nil {
			return 0, fmt.Errorf("error creating WebP file: %v", err)
		}
		defer file.Close()
		qualityFloat := float32(quality)
		options, err := encoder.NewLossyEncoderOptions(encoder.PresetPhoto, qualityFloat)
		if err != nil {
			return 0, fmt.Errorf("error creating encoder options: %v", err)
		}
		enc, err := encoder.NewEncoder(img, options)
		if err != nil {
			return 0, fmt.Errorf("error creating encoder: %v", err)
		}
		if err = enc.Encode(file); err != nil {
			return 0, fmt.Errorf("error encoding WebP: %v", err)
		}
	} else if format == "png" {
		file, err = os.Create(path)
		if err != nil {
			return 0, fmt.Errorf("error creating PNG file: %v", err)
		}
		defer file.Close()
		if err := png.Encode(file, img); err != nil {
			return 0, fmt.Errorf("error encoding PNG: %v", err)
		}
	} else {
		return 0, fmt.Errorf("unsupported format")
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("error getting file info: %v", err)
	}
	return int(fileInfo.Size() / 1024), nil
}

func OptimizeQuality(img image.Image, outputFormat, base string, targetSizeKB int) (int, error) {
	if outputFormat == "png" {
		return 80, nil
	}

	low := 1
	high := 100
	bestQuality := 0

	for low <= high {
		mid := low + (high-low)/2
		path := base + "_" + fmt.Sprintf("%d", mid) + "." + outputFormat
		size, err := saveAndGetSize(img, mid, outputFormat, path)
		if err != nil {
			os.Remove(path)
			return 0, err
		}
		fmt.Printf("Testing quality %d, size %d KB\n", mid, size)
		if size <= targetSizeKB {
			bestQuality = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
		os.Remove(path)
	}

	if bestQuality == 0 {
		return 1, fmt.Errorf("could not optimize quality for %s to fit %d KB", outputFormat, targetSizeKB)
	}
	return bestQuality, nil
}
