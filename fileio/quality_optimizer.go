package fileio

import (
	"fmt"
	"image"
	"image/jpeg"
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
		qualityFloat := float32(quality) / 100.0
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
	var low, high, bestQuality int
	low, high = 1, 100
	bestQuality = 20
	for low <= high {
		mid := (low + high) / 2
		size, err := saveAndGetSize(img, mid, outputFormat, base+"_"+fmt.Sprintf("%d", mid)+"."+outputFormat)
		if err != nil {
			return 0, err
		}
		fmt.Printf("Testing quality %d, size %d KB\n", mid, size)
		if size <= targetSizeKB {
			bestQuality = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
		os.Remove(base + "_" + fmt.Sprintf("%d", mid) + "." + outputFormat)
	}
	for i := bestQuality; i <= 100; i++ {
		size, err := saveAndGetSize(img, i, outputFormat, base+"_"+fmt.Sprintf("%d", i)+"."+outputFormat)
		if err != nil {
			return 0, err
		}
		if size <= targetSizeKB {
			bestQuality = i
		} else {
			break
		}
		os.Remove(base + "_" + fmt.Sprintf("%d", i) + "." + outputFormat)
	}
	for i := 1; i <= 100; i++ {
		os.Remove(base + "_" + fmt.Sprintf("%d", i) + "." + outputFormat)
	}
	return bestQuality, nil
}
