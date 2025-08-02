package main

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

func main() {
	files, err := os.ReadDir("Images")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
	fmt.Printf("Found %d files in Images/: %v\n", len(files), listFiles(files))
	err = os.Mkdir("Images_watermarked", 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println("Error creating directory:", err)
		return
	}
	watermark, err := imaging.Open("Images/watermark.png")
	if err != nil {
		fmt.Println("Error loading watermark:", err)
		return
	}
	for _, file := range files {
		fmt.Printf("Checking file: %s\n", file.Name())
		if file.Name() == "watermark.png" {
			fmt.Println("Skipping watermark.png")
			continue
		}
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".webp" {
			fmt.Printf("Processing image: %s\n", file.Name())
			img, err := imaging.Open("Images/" + file.Name())
			if err != nil {
				fmt.Printf("Error loading image %s: %v\n", file.Name(), err)
				continue
			}
			width := img.Bounds().Dx()
			height := img.Bounds().Dy()
			if width > 1200 || height > 1200 {
				img = imaging.Fit(img, 1200, 1200, imaging.Lanczos)
				width = img.Bounds().Dx()
				height = img.Bounds().Dy()
				fmt.Printf("Resized image %s to %dx%d\n", file.Name(), width, height)
			}
			watermarkResized := imaging.Resize(watermark, width, height, imaging.Lanczos)
			bounds := watermarkResized.Bounds()
			transparentWatermark := image.NewNRGBA(bounds)
			draw.Draw(transparentWatermark, bounds, watermarkResized, image.Point{0, 0}, draw.Src)
			// alpha layout cfg
			/*
			   for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			       for x := bounds.Min.X; x < bounds.Max.X; x++ {
			           r, g, b, a := transparentWatermark.At(x, y).RGBA()
			           transparentWatermark.SetNRGBA(x, y, color.NRGBA{
			               R: uint8(r >> 8),
			               G: uint8(g >> 8),
			               B: uint8(b >> 8),
			               A: uint8((a >> 8) * 75 / 100), // 75% от оригинального альфа-слоя
			           })
			       }
			   }
			*/
			result := image.NewNRGBA(img.Bounds())
			draw.Draw(result, img.Bounds(), img, image.Point{0, 0}, draw.Src)
			draw.Draw(result, transparentWatermark.Bounds(), transparentWatermark, image.Point{0, 0}, draw.Over)
			base := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			outputPath := "Images_watermarked/" + base + ".png"
			err = imaging.Save(result, outputPath)
			if err != nil {
				fmt.Printf("Error saving image %s: %v\n", outputPath, err)
				continue
			}
			fmt.Printf("Processed image: %s\n", outputPath)
		} else {
			fmt.Printf("Skipping file %s: unsupported extension %s\n", file.Name(), ext)
		}
	}
}

func listFiles(files []os.DirEntry) []string {
	names := make([]string, len(files))
	for i, file := range files {
		names[i] = file.Name()
	}
	return names
}
