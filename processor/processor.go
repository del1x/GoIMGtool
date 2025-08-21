package processor

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/del1x/GoIMGtool/config"
	"github.com/del1x/GoIMGtool/fileio"
	"github.com/disintegration/imaging"
)

type ProgressCallback func(current, total int, img *canvas.Image, fileName string)

type FileHandler interface {
	LoadImage(path string) (image.Image, error)
	SaveImage(img image.Image, path, format string, cfg *config.Config) error
	CreateDir(path string) error
	ReadDir(path string) ([]os.DirEntry, error)
}

type ImageProcessor struct {
	Watermark     image.Image
	OutputDir     string
	Config        *config.Config
	WatermarkMode string
	FileHandler   FileHandler
}

func NewImageProcessor(watermarkPath string, cfg *config.Config, fileHandler FileHandler) (*ImageProcessor, error) {
	watermark, err := fileHandler.LoadImage(watermarkPath)
	if err != nil {
		return nil, fmt.Errorf("error loading watermark: %v", err)
	}
	if watermark == nil {
		return nil, fmt.Errorf("watermark image is nil")
	}
	return &ImageProcessor{
		Watermark:     watermark,
		OutputDir:     "Images_watermarked",
		Config:        cfg,
		WatermarkMode: "crop",
		FileHandler:   fileHandler,
	}, nil
}

func (p *ImageProcessor) ProcessFolder(imageDir, outputFormat string, progress ProgressCallback) error {
	startTime := time.Now()
	files, err := p.FileHandler.ReadDir(imageDir)
	if err != nil {
		return fmt.Errorf("error reading directory: %v", err)
	}
	total := len(files)

	if err := p.setupOutputDir(); err != nil {
		return err
	}

	const maxGoroutines = 4
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxGoroutines)
	current := 0

	for _, file := range files {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(f os.DirEntry) {
			defer wg.Done()
			defer func() { <-semaphore }()

			outputPath, err := p.processFile(imageDir, f, outputFormat)
			if err != nil {
				fmt.Printf("Error processing file %s: %v\n", f.Name(), err)
				return
			}
			current++
			if progress != nil {
				fyne.Do(func() {
					canvasImg := p.displayImageInUI(outputPath)
					progress(current, total, canvasImg, f.Name())
				})
			}
		}(file)
	}

	wg.Wait()
	if progress != nil {
		fyne.Do(func() {
			progress(total, total, nil, "")
		})
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Processing completed in %v\n", elapsed)
	return nil
}

func (p *ImageProcessor) processFile(imageDir string, file os.DirEntry, outputFormat string) (string, error) {
	if p.isWatermarkFile(file, imageDir) {
		fmt.Println("Skipping watermark.png")
		return "", nil
	}
	if !p.isSupportedExtension(file) {
		fmt.Printf("Skipping file %s: unsupported extension %s\n", file.Name(), filepath.Ext(file.Name()))
		return "", nil
	}

	img, err := p.FileHandler.LoadImage(filepath.Join(imageDir, file.Name()))
	if err != nil {
		return "", fmt.Errorf("failed to load image %s: %v", file.Name(), err)
	}
	if img == nil {
		return "", fmt.Errorf("loaded image for %s is nil", file.Name())
	}
	img, err = p.resizeImage(img)
	if err != nil {
		return "", err
	}
	result, err := p.applyWatermark(img)
	if err != nil {
		return "", err
	}
	outputPath := filepath.Join(p.OutputDir, strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))+"."+outputFormat)
	if err := p.FileHandler.SaveImage(result, outputPath, outputFormat, p.Config); err != nil {
		return "", err
	}
	fmt.Printf("Image saved to %s\n", outputPath)
	return outputPath, nil
}

func (p *ImageProcessor) resizeImage(img image.Image) (image.Image, error) {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	if width > p.Config.MaxWidth || height > p.Config.MaxHeight {
		img = fileio.HandleImageResize(img, p.Config)
		fmt.Printf("Resized image to %dx%d\n", img.Bounds().Dx(), img.Bounds().Dy())
	}
	return img, nil
}

func (p *ImageProcessor) applyWatermark(img image.Image) (image.Image, error) {
	watermark := p.prepareWatermark(img)
	bounds := img.Bounds()
	transparentWatermark := image.NewNRGBA(bounds)
	draw.Draw(transparentWatermark, bounds, watermark, image.Point{0, 0}, draw.Src)

	result := image.NewNRGBA(bounds)
	draw.Draw(result, bounds, img, image.Point{0, 0}, draw.Src)
	draw.Draw(result, bounds, transparentWatermark, image.Point{0, 0}, draw.Over)
	return result, nil
}

func (p *ImageProcessor) setupOutputDir() error {
	return p.FileHandler.CreateDir(p.OutputDir)
}

func (p *ImageProcessor) isWatermarkFile(file os.DirEntry, imageDir string) bool {
	return file.Name() == filepath.Base(filepath.Join(imageDir, "watermark.png"))
}

func (p *ImageProcessor) isSupportedExtension(file os.DirEntry) bool {
	ext := strings.ToLower(filepath.Ext(file.Name()))
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".webp"
}

func (p *ImageProcessor) prepareWatermark(img image.Image) image.Image {
	wmBounds := p.Watermark.Bounds()
	imgBounds := img.Bounds()

	if p.WatermarkMode == "resize" {
		return imaging.Resize(p.Watermark, imgBounds.Dx(), imgBounds.Dy(), imaging.Lanczos)
	}

	if wmBounds.Dx() <= imgBounds.Dx() && wmBounds.Dy() <= imgBounds.Dy() {
		return imaging.Resize(p.Watermark, imgBounds.Dx(), imgBounds.Dy(), imaging.Lanczos)
	}

	cropX := (wmBounds.Dx() - imgBounds.Dx()) / 2
	cropY := (wmBounds.Dy() - imgBounds.Dy()) / 2
	cropRect := image.Rect(cropX, cropY, cropX+imgBounds.Dx(), cropY+imgBounds.Dy())

	croppedWatermark := imaging.Crop(p.Watermark, cropRect)
	fmt.Printf("Cropped watermark to %dx%d from center\n", croppedWatermark.Bounds().Dx(), croppedWatermark.Bounds().Dy())

	return croppedWatermark
}

func (p *ImageProcessor) displayImageInUI(path string) *canvas.Image {
	if path == "" {
		return nil
	}
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("Error decoding image for UI: %v", err)
		return nil
	}
	canvasImg := canvas.NewImageFromImage(img)
	canvasImg.FillMode = canvas.ImageFillContain
	canvasImg.SetMinSize(fyne.NewSize(200, 200))
	return canvasImg
}
