package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"image/png"
)

func main() {
	if err := os.Mkdir("test_images", 0755); err != nil && !os.IsExist(err) {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	for i := 1; i <= 50; i++ {
		img := image.NewRGBA(image.Rect(0, 0, 1200, 1200))
		gray := color.RGBA{128, 128, 128, 255}
		for y := 0; y < 1200; y++ {
			for x := 0; x < 1200; x++ {
				img.Set(x, y, gray)
			}
		}

		filename := fmt.Sprintf("test_images/test_%d.png", i)
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Error creating file %s: %v\n", filename, err)
			continue
		}
		if err := png.Encode(file, img); err != nil {
			fmt.Printf("Error encoding file %s: %v\n", filename, err)
			file.Close()
			continue
		}
		file.Close()
		fmt.Printf("Generated %s\n", filename)
	}
	fmt.Println("Generation complete!")
}
