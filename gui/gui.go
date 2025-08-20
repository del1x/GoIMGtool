package gui

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/del1x/GoIMGtool/config"
	"github.com/del1x/GoIMGtool/fileio"
	"github.com/del1x/GoIMGtool/processor"
)

type fileioHandler struct{}

func (f *fileioHandler) LoadImage(path string) (image.Image, error) {
	return fileio.LoadImage(path)
}

func (f *fileioHandler) SaveImage(img image.Image, path, format string, cfg *config.Config) error {
	return fileio.NewImageProcessor().SaveImage(img, path, format, cfg)
}

func (f *fileioHandler) CreateDir(path string) error {
	return fileio.CreateDir(path)
}

func (f *fileioHandler) ReadDir(path string) ([]os.DirEntry, error) {
	return fileio.ReadDir(path)
}

func SetupGUI(window fyne.Window) {
	cfg := config.JpgConfig()
	fileHandler := &fileioHandler{}
	watermarkMode := "crop"

	widthLabel := widget.NewLabel("Max Width (100-4096):")
	heightLabel := widget.NewLabel("Max Height (100-4096):")

	watermarkLabel := widget.NewLabel("Watermark file:")
	watermarkEntry := widget.NewEntry()
	watermarkEntry.SetPlaceHolder("Select watermark.png")

	imageDirLabel := widget.NewLabel("Image folder:")
	imageDirEntry := widget.NewEntry()
	imageDirEntry.SetPlaceHolder("Select image folder")

	formatLabel := widget.NewLabel("Output format:")
	formatSelect := widget.NewSelect([]string{"jpg", "webp", "png"}, func(s string) {
		cfg.OutputFormat = s
	})
	formatSelect.SetSelected("jpg")

	qualityLabel := widget.NewLabel("Quality (1-100):")
	qualityEntry := widget.NewEntry()
	qualityEntry.SetText(strconv.Itoa(cfg.Quality))
	qualityEntry.OnChanged = func(s string) {
		if q, err := strconv.Atoi(s); err == nil && q >= 1 && q <= 100 {
			cfg.Quality = q
		} else {
			qualityEntry.SetText(strconv.Itoa(cfg.Quality))
		}
	}

	webSizeHint := widget.NewLabel("Note: For web, target size ≤100 KB is optimal")

	widthEntry := widget.NewEntry()
	widthEntry.SetText(strconv.Itoa(cfg.MaxWidth))
	widthEntry.OnChanged = func(s string) {
		if w, err := strconv.Atoi(s); err == nil && w >= 1 {
			if watermarkEntry.Text != "" {
				if img, err := fileHandler.LoadImage(watermarkEntry.Text); err == nil && img != nil {
					if w > img.Bounds().Dx() {
						dialog.ShowInformation("Warning", fmt.Sprintf("Width exceeds watermark width (%d px)", img.Bounds().Dx()), window)
						widthEntry.SetText(strconv.Itoa(img.Bounds().Dx()))
						cfg.MaxWidth = img.Bounds().Dx()
						return
					}
				}
			}
			cfg.MaxWidth = w
		} else {
			widthEntry.SetText(strconv.Itoa(cfg.MaxWidth))
		}
	}

	heightEntry := widget.NewEntry()
	heightEntry.SetText(strconv.Itoa(cfg.MaxHeight))
	heightEntry.OnChanged = func(s string) {
		if h, err := strconv.Atoi(s); err == nil && h >= 1 {
			if watermarkEntry.Text != "" {
				if img, err := fileHandler.LoadImage(watermarkEntry.Text); err == nil && img != nil {
					if h > img.Bounds().Dy() {
						dialog.ShowInformation("Warning", fmt.Sprintf("Height exceeds watermark height (%d px)", img.Bounds().Dy()), window)
						heightEntry.SetText(strconv.Itoa(img.Bounds().Dy()))
						cfg.MaxHeight = img.Bounds().Dy()
						return
					}
				}
			}
			cfg.MaxHeight = h
		} else {
			heightEntry.SetText(strconv.Itoa(cfg.MaxHeight))
		}
	}

	targetSizeLabel := widget.NewLabel("Target Size (KB, 50-5000):")
	targetSizeEntry := widget.NewEntry()
	targetSizeEntry.SetText("100")
	targetSizeEntry.OnChanged = func(s string) {
		if size, err := strconv.Atoi(s); err == nil && size >= 50 && size <= 5000 {
			fileio.NewImageProcessor().TargetSizeKB = size
		} else {
			targetSizeEntry.SetText("100")
		}
	}

	watermarkModeLabel := widget.NewLabel("Watermark Mode:")
	watermarkModeSelect := widget.NewSelect([]string{"crop", "resize"}, func(s string) {
		watermarkMode = s
	})
	watermarkModeSelect.SetSelected("crop")

	currentFileLabel := widget.NewLabel("Processing: None")
	progress := widget.NewProgressBar()
	imageContainer := container.NewVBox()
	scrollContainer := container.NewVScroll(imageContainer)
	scrollContainer.SetMinSize(fyne.NewSize(400, 200))

	fileButton := widget.NewButton("Browse...", func() {
		dialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				dialog.ShowInformation("Error", "Failed to select watermark file!", window)
				return
			}
			path := reader.URI().Path()
			if _, err := os.Stat(path); err != nil {
				dialog.ShowInformation("Error", fmt.Sprintf("Invalid file: %v", err), window)
				return
			}
			watermarkEntry.SetText(path)
			if img, err := fileHandler.LoadImage(path); err == nil && img != nil {
				maxWidth := img.Bounds().Dx()
				maxHeight := img.Bounds().Dy()
				widthLabel.SetText(fmt.Sprintf("Max Width (100-%d):", maxWidth))
				heightLabel.SetText(fmt.Sprintf("Max Height (100-%d):", maxHeight))
				if cfg.MaxWidth > maxWidth {
					widthEntry.SetText(strconv.Itoa(maxWidth))
					cfg.MaxWidth = maxWidth
				}
				if cfg.MaxHeight > maxHeight {
					heightEntry.SetText(strconv.Itoa(maxHeight))
					cfg.MaxHeight = maxHeight
				}
			} else {
				widthLabel.SetText("Max Width (100-4096):")
				heightLabel.SetText("Max Height (100-4096):")
			}
		}, window)
		dialog.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg", ".webp"}))
		dialog.Show()
	})

	folderButton := widget.NewButton("Browse Folder...", func() {
		dialog := dialog.NewFolderOpen(func(reader fyne.ListableURI, err error) {
			if err != nil || reader == nil {
				dialog.ShowInformation("Error", "Failed to select folder!", window)
				return
			}
			path := reader.Path()
			if _, err := os.Stat(path); err != nil {
				dialog.ShowInformation("Error", fmt.Sprintf("Invalid folder: %v", err), window)
				return
			}
			imageDirEntry.SetText(path)
		}, window)
		dialog.Show()
	})

	processButton := widget.NewButton("Process", func() {
		if watermarkEntry.Text == "" || imageDirEntry.Text == "" {
			dialog.ShowInformation("Error", "Please select a watermark file and image folder!", window)
			return
		}

		if size, err := strconv.Atoi(targetSizeEntry.Text); err == nil && size >= 50 && size <= 5000 {
			fileio.NewImageProcessor().TargetSizeKB = size
		} else {
			targetSizeEntry.SetText("100")
			fileio.NewImageProcessor().TargetSizeKB = 100
		}

		progress.SetValue(0)
		imageContainer.Objects = nil
		currentFileLabel.SetText("Processing: None")
		processor, err := processor.NewImageProcessor(watermarkEntry.Text, cfg, fileHandler)
		if err != nil {
			dialog.ShowInformation("Error", fmt.Sprintf("Failed to initialize processor: %v", err), window)
			return
		}
		processor.WatermarkMode = watermarkMode
		err = processor.ProcessFolder(imageDirEntry.Text, cfg.OutputFormat, func(current, total int, img *canvas.Image, fileName string) {
			fyne.Do(func() {
				if total > 0 {
					progress.SetValue(float64(current) / float64(total))
					if img != nil {
						currentFileLabel.SetText(fmt.Sprintf("Processing: %s", fileName))
						fileInfo, err := os.Stat(filepath.Join("Images_watermarked", fileName+"."+cfg.OutputFormat))
						sizeKB := "N/A"
						if err == nil {
							sizeKB = fmt.Sprintf("%d KB", fileInfo.Size()/1024)
						}
						imageContainer.Add(container.NewVBox(
							widget.NewLabel(fmt.Sprintf("%s (%s)", fileName, sizeKB)),
							img,
						))
						scrollContainer.Refresh()
					}
					window.Canvas().Refresh(progress)
					window.Canvas().Refresh(currentFileLabel)
				}
			})
		})
		if err != nil {
			dialog.ShowInformation("Error", fmt.Sprintf("Processing failed: %v", err), window)
		} else {
			currentFileLabel.SetText("Processing: Done")
		}
	})

	content := container.NewVBox(
		watermarkLabel, watermarkEntry, fileButton,
		imageDirLabel, imageDirEntry, folderButton,
		formatLabel, formatSelect,
		qualityLabel, qualityEntry,
		webSizeHint,
		widthLabel, widthEntry,
		heightLabel, heightEntry,
		targetSizeLabel, targetSizeEntry,
		watermarkModeLabel, watermarkModeSelect,
		processButton, currentFileLabel, progress, scrollContainer,
	)

	window.SetContent(content)
	window.Resize(fyne.NewSize(400, 700))
}
