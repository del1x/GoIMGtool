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
	"github.com/del1x/GoIMGtool/config"
	"github.com/disintegration/imaging"
)

type ProgressCallback func(current, total int)

type FileHandler interface {
	LoadImage(path string) (image.Image, error)
	SaveImage(img image.Image, path, format string, cfg *config.Config) error
	CreateDir(path string) error
	ReadDir(path string) ([]os.DirEntry, error)
}

type ImageProcessor struct {
	Watermark   image.Image
	OutputDir   string
	Config      *config.Config
	FileHandler FileHandler
}

func NewImageProcessor(watermarkPath string, cfg *config.Config, fileHandler FileHandler) (*ImageProcessor, error) {
	watermark, err := fileHandler.LoadImage(watermarkPath)
	if err != nil {
		return nil, fmt.Errorf("error loading watermark: %v", err)
	}
	return &ImageProcessor{
		Watermark:   watermark,
		OutputDir:   "Images_watermarked",
		Config:      cfg,
		FileHandler: fileHandler,
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

			if err := p.processFile(imageDir, f, outputFormat); err != nil {
				fmt.Printf("Error processing file %s: %v\n", f.Name(), err)
			}
			current++
			if progress != nil {
				fyne.Do(func() {
					progress(current, total)
				})
			}
		}(file)
	}

	wg.Wait()
	if progress != nil {
		fyne.Do(func() {
			progress(total, total)
		})
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Processing completed in %v\n", elapsed)
	return nil
}

func (p *ImageProcessor) processFile(imageDir string, file os.DirEntry, outputFormat string) error {
	if p.isWatermarkFile(file, imageDir) {
		fmt.Println("Skipping watermark.png")
		return nil
	}
	if !p.isSupportedExtension(file) {
		fmt.Printf("Skipping file %s: unsupported extension %s\n", file.Name(), filepath.Ext(file.Name()))
		return nil
	}

	img, err := p.FileHandler.LoadImage(filepath.Join(imageDir, file.Name()))
	if err != nil {
		return err
	}
	img, err = p.resizeImage(img)
	if err != nil {
		return err
	}
	result, err := p.applyWatermark(img)
	if err != nil {
		return err
	}
	return p.FileHandler.SaveImage(result, filepath.Join(p.OutputDir, file.Name()), outputFormat, p.Config)
}

func (p *ImageProcessor) resizeImage(img image.Image) (image.Image, error) {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	if width > 1200 || height > 1200 {
		img = imaging.Fit(img, 1200, 1200, imaging.Lanczos)
		fmt.Printf("Resized image to %dx%d\n", img.Bounds().Dx(), img.Bounds().Dy())
	}
	return img, nil
}

func (p *ImageProcessor) applyWatermark(img image.Image) (image.Image, error) {
	watermark := p.prepareWatermark(img)
	bounds := watermark.Bounds()
	transparentWatermark := image.NewNRGBA(bounds)
	draw.Draw(transparentWatermark, bounds, watermark, image.Point{0, 0}, draw.Src)

	result := image.NewNRGBA(img.Bounds())
	draw.Draw(result, img.Bounds(), img, image.Point{0, 0}, draw.Src)
	draw.Draw(result, transparentWatermark.Bounds(), transparentWatermark, image.Point{0, 0}, draw.Over)
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
	if p.Watermark.Bounds().Dx() > img.Bounds().Dx() || p.Watermark.Bounds().Dy() > img.Bounds().Dy() {
		return imaging.CropCenter(p.Watermark, img.Bounds().Dx(), img.Bounds().Dy())
	}
	return imaging.Resize(p.Watermark, img.Bounds().Dx(), img.Bounds().Dy(), imaging.Lanczos)
}
