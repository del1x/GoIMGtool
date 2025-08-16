package fileio

import (
	"fmt"
	"image"
	"os"

	"github.com/disintegration/imaging"
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
