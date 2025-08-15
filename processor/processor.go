package processor

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"path/filepath"
	"strings"

	"github.com/del1x/GoIMGtool/config"
	"github.com/del1x/GoIMGtool/fileio"
	"github.com/disintegration/imaging"
)

type ProgressCallback func(current, total int)

type ImageProcessor struct {
	Watermark image.Image
	OutputDir string
	Config    *config.Config
}

func NewImageProcessor(watermarkPath string, cfg *config.Config) (*ImageProcessor, error) {
	watermark, err := fileio.LoadImage(watermarkPath)
	if err != nil {
		return nil, fmt.Errorf("error loading watermark: %v", err)
	}
	return &ImageProcessor{
		Watermark: watermark,
		OutputDir: "Images_watermarked",
		Config:    cfg,
	}, nil
}

func (p *ImageProcessor) ProcessFolder(imageDir, outputFormat string, progress ProgressCallback) error {
	files, err := fileio.ReadDir(imageDir)
	if err != nil {
		return fmt.Errorf("error reading directory: %v", err)
	}
	total := len(files)
	if progress != nil {
		progress(0, total)
	}
	if err := fileio.CreateDir(p.OutputDir); err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
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
	img, err := fileio.LoadImage(filepath.Join(imageDir, file.Name()))
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
	return fileio.SaveImage(result, filepath.Join(p.OutputDir, file.Name()), outputFormat)
}

func resizeImage(img image.Image) (image.Image, error) {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	if width > 1200 || height > 1200 {
		img = imaging.Fit(img, 1200, 1200, imaging.Lanczos)
		fmt.Printf("Resized image to %dx%d\n", img.Bounds().Dx(), img.Bounds().Dy())
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
