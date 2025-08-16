package fileio

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
)

func ManageExif(img image.Image, outputPath string, quality int) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating JPEG file: %v", err)
	}
	defer file.Close()

	if err := jpeg.Encode(file, img, &jpeg.Options{Quality: quality}); err != nil {
		return fmt.Errorf("error encoding JPEG: %v", err)
	}

	fmt.Printf("Encoding JPEG with optimized quality: %d\n", quality)
	return nil
}
