package gui

import (
	"image"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
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

func SetupGUI(w fyne.Window) {
	cfg := config.DefaultConfig()

	fileHandler := &fileioHandler{}

	watermarkLabel := widget.NewLabel("Watermark file:")
	watermarkEntry := widget.NewEntry()
	watermarkEntry.SetPlaceHolder("Select watermark.png")

	imageDirLabel := widget.NewLabel("Image folder:")
	imageDirEntry := widget.NewEntry()
	imageDirEntry.SetPlaceHolder("Select image folder")

	formatLabel := widget.NewLabel("Output format:")
	formatSelect := widget.NewSelect([]string{"png", "webp", "jpg"}, func(s string) {
		cfg.OutputFormat = s
	})
	formatSelect.SetSelected("png")

	progress := widget.NewProgressBar()

	fileButton := widget.NewButton("Browse...", func() {
		dialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			watermarkEntry.SetText(reader.URI().Path())
		}, w)
		dialog.SetFilter(storage.NewExtensionFileFilter([]string{".png"}))
		dialog.Show()
	})

	folderButton := widget.NewButton("Browse Folder...", func() {
		dialog := dialog.NewFolderOpen(func(reader fyne.ListableURI, err error) {
			if err != nil || reader == nil {
				return
			}
			imageDirEntry.SetText(reader.Path())
		}, w)
		dialog.Show()
	})

	qualityEntry := widget.NewEntry()
	qualityEntry.SetText(strconv.Itoa(cfg.Quality))
	qualityEntry.OnChanged = func(s string) {
		if q, err := strconv.Atoi(s); err == nil {
			if q >= 1 && q <= 100 {
				cfg.Quality = q
				qualityEntry.SetText(strconv.Itoa(q))
			} else {
				qualityEntry.SetText(strconv.Itoa(cfg.Quality))
			}
		} else {
			qualityEntry.SetText(strconv.Itoa(cfg.Quality))
		}
	}

	processButton := widget.NewButton("Process", func() {
		if watermarkEntry.Text == "" || imageDirEntry.Text == "" {
			dialog.ShowInformation("Error", "Please select a watermark file and image folder!", w)
			return
		}
		progress.SetValue(0)
		processor, err := processor.NewImageProcessor(watermarkEntry.Text, cfg, fileHandler)
		if err != nil {
			dialog.ShowInformation("Error", err.Error(), w)
			return
		}
		err = processor.ProcessFolder(imageDirEntry.Text, cfg.OutputFormat, func(current, total int) {
			fyne.Do(func() {
				if total > 0 {
					progress.SetValue(float64(current) / float64(total))
					w.Canvas().Refresh(progress)
				}
			})
		})
		if err != nil {
			dialog.ShowInformation("Error", err.Error(), w)
		}
	})

	content := container.NewVBox(
		watermarkLabel, watermarkEntry, fileButton,
		imageDirLabel, imageDirEntry, folderButton,
		formatLabel, formatSelect,
		widget.NewLabel("Quality (1-100):"), qualityEntry,
		processButton, progress,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 350))
}
