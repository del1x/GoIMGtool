package fileio

import (
	"bytes"
	"fmt"
	"image"
)

func saveAndGetSize(img image.Image, quality int, format string) (int, error) {
	var buf bytes.Buffer
	encoder, err := GetEncoder(format)
	if err != nil {
		return 0, err
	}
	if err := encoder.Encode(img, &buf, quality); err != nil {
		return 0, fmt.Errorf("error encoding %s: %v", format, err)
	}
	return buf.Len() / 1024, nil
}

func OptimizeQuality(img image.Image, outputFormat, base string, targetSizeKB int) (int, error) {
	if outputFormat == "png" {
		return 80, nil
	}

	low, high, bestQuality := 1, 100, 0
	for low <= high {
		mid := low + (high-low)/2
		size, err := saveAndGetSize(img, mid, outputFormat)
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
	}

	if bestQuality == 0 {
		return 1, fmt.Errorf("could not optimize quality for %s to fit %d KB", outputFormat, targetSizeKB)
	}
	return bestQuality, nil
}
