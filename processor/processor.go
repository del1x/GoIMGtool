package processor

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/kolesa-team/go-webp/encoder"
)

type ProgressCallback func(current, total int)

type ImageProcessor struct {
	Watermark image.Image
	OutputDir string
}

func NewImageProcessor(watermarkPath string) (*ImageProcessor, error) {
	watermark, err := imaging.Open(watermarkPath)
	if err != nil {
		return nil, fmt.Errorf("error loading watermark: %v", err)
	}
	return &ImageProcessor{Watermark: watermark, OutputDir: "Images_watermarked"}, nil
}

func (p *ImageProcessor) ProcessFolder(imageDir string, outputFormat string, progress ProgressCallback) error {
	files, err := os.ReadDir(imageDir)
	if err != nil {
		return fmt.Errorf("error reading directory: %v", err)
	}
	total := len(files)
	if progress != nil {
		progress(0, total)
	}
	fmt.Printf("Found %d files in %s/: %v\n", total, imageDir, listFiles(files))
	err = os.Mkdir(p.OutputDir, 0755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("error creating directory: %v", err)
	}
	for i, file := range files {
		if progress != nil {
			progress(i+1, total)
		}
		if err := p.processFile(imageDir, file, outputFormat); err != nil {
			fmt.Printf("Error processing file %s: %v\n", file.Name(), err)
			continue
		}
	}
	if progress != nil {
		progress(total, total)
	}
	return nil
}

func (p *ImageProcessor) processFile(imageDir string, file os.DirEntry, outputFormat string) error {
	if file.Name() == filepath.Base(filepath.Join(imageDir, "watermark.png")) {
		fmt.Println("Skipping watermark.png")
		return nil
	}
	ext := strings.ToLower(filepath.Ext(file.Name()))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".webp" {
		fmt.Printf("Skipping file %s: unsupported extension %s\n", file.Name(), ext)
		return nil
	}
	fmt.Printf("Processing image: %s\n", file.Name())
	img, err := loadImage(filepath.Join(imageDir, file.Name()))
	if err != nil {
		return err
	}
	img, err = resizeImage(img)
	if err != nil {
		return err
	}
	result, err := applyWatermark(img, p.Watermark)
	if err != nil {
		return err
	}
	return saveImage(result, file.Name(), outputFormat)
}

func loadImage(path string) (image.Image, error) {
	img, err := imaging.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error loading image: %v", err)
	}
	return img, nil
}

func resizeImage(img image.Image) (image.Image, error) {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	if width > 1200 || height > 1200 {
		img = imaging.Fit(img, 1200, 1200, imaging.Lanczos)
		width = img.Bounds().Dx()
		height = img.Bounds().Dy()
		fmt.Printf("Resized image to %dx%d\n", width, height)
	}
	return img, nil
}

func applyWatermark(img, watermark image.Image) (image.Image, error) {
	watermarkResized := imaging.Resize(watermark, img.Bounds().Dx(), img.Bounds().Dy(), imaging.Lanczos)
	bounds := watermarkResized.Bounds()
	transparentWatermark := image.NewNRGBA(bounds)
	draw.Draw(transparentWatermark, bounds, watermarkResized, image.Point{0, 0}, draw.Src)
	result := image.NewNRGBA(img.Bounds())
	draw.Draw(result, img.Bounds(), img, image.Point{0, 0}, draw.Src)
	draw.Draw(result, transparentWatermark.Bounds(), transparentWatermark, image.Point{0, 0}, draw.Over)
	return result, nil
}

func saveImage(img image.Image, filename, outputFormat string) error {
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	var outputPath string
	if outputFormat == "webp" {
		outputPath = filepath.Join("Images_watermarked", base+".webp")
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
		outputPath = filepath.Join("Images_watermarked", base+".png")
		err := imaging.Save(img, outputPath)
		if err != nil {
			return fmt.Errorf("error saving image: %v", err)
		}
	}
	fmt.Printf("Processed image: %s\n", outputPath)
	return nil
}

func listFiles(files []os.DirEntry) []string {
	names := make([]string, len(files))
	for i, file := range files {
		names[i] = file.Name()
	}
	return names
}
