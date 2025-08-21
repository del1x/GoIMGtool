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

type FileHandler interface {
	LoadImage(path string) (image.Image, error)
	SaveImage(img image.Image, path, format string, cfg *config.Config) error
	CreateDir(path string) error
	ReadDir(path string) ([]os.DirEntry, error)
}

type GUI struct {
	window        fyne.Window
	cfg           *config.Config
	fileHandler   FileHandler
	watermarkMode string
	currentLocale string
	components    *GUIComponents
}

type GUIComponents struct {
	widthLabel, heightLabel, watermarkLabel, imageDirLabel, formatLabel, qualityLabel, targetSizeLabel, watermarkModeLabel, languageLabel, webSizeHintLabel *widget.Label
	widthEntry, heightEntry, qualityEntry, targetSizeEntry, watermarkEntry, imageDirEntry                                                                   *widget.Entry
	formatSelect, watermarkModeSelect, languageSelect                                                                                                       *widget.Select
	currentFileLabel                                                                                                                                        *widget.Label
	progress                                                                                                                                                *widget.ProgressBar
	imageContainer                                                                                                                                          *fyne.Container
	scrollContainer                                                                                                                                         *container.Scroll
	fileButton, folderButton                                                                                                                                *widget.Button
}

func NewGUI(window fyne.Window, cfg *config.Config, fileHandler FileHandler) *GUI {
	return &GUI{
		window:        window,
		cfg:           cfg,
		fileHandler:   fileHandler,
		watermarkMode: "crop",
		currentLocale: "en",
		components:    &GUIComponents{},
	}
}

func (g *GUI) Setup() {
	g.initComponents()
	content := container.NewVBox(
		g.components.languageLabel, g.components.languageSelect,
		g.components.watermarkLabel, g.components.watermarkEntry, g.components.fileButton,
		g.components.imageDirLabel, g.components.imageDirEntry, g.components.folderButton,
		g.components.formatLabel, g.components.formatSelect,
		g.components.qualityLabel, g.components.qualityEntry,
		g.components.webSizeHintLabel,
		g.components.widthLabel, g.components.widthEntry,
		g.components.heightLabel, g.components.heightEntry,
		g.components.targetSizeLabel, g.components.targetSizeEntry,
		g.components.watermarkModeLabel, g.components.watermarkModeSelect,
		g.createProcessButton(),
		g.components.currentFileLabel,
		g.components.progress,
		g.components.scrollContainer,
	)
	g.window.SetContent(content)
	g.window.Resize(fyne.NewSize(400, 700))
}

func (g *GUI) initComponents() {
	g.components.languageLabel = widget.NewLabel(locales[g.currentLocale].LanguageLabel)
	g.components.widthLabel = widget.NewLabel(locales[g.currentLocale].WidthLabel)
	g.components.heightLabel = widget.NewLabel(locales[g.currentLocale].HeightLabel)
	g.components.watermarkLabel = widget.NewLabel(locales[g.currentLocale].WatermarkLabel)
	g.components.imageDirLabel = widget.NewLabel(locales[g.currentLocale].ImageDirLabel)
	g.components.formatLabel = widget.NewLabel(locales[g.currentLocale].FormatLabel)
	g.components.qualityLabel = widget.NewLabel(locales[g.currentLocale].QualityLabel)
	g.components.targetSizeLabel = widget.NewLabel(locales[g.currentLocale].TargetSizeLabel)
	g.components.watermarkModeLabel = widget.NewLabel(locales[g.currentLocale].WatermarkModeLabel)
	g.components.currentFileLabel = widget.NewLabel(locales[g.currentLocale].CurrentFileLabel)
	g.components.webSizeHintLabel = widget.NewLabel(locales[g.currentLocale].WebSizeHint)

	g.components.watermarkEntry = widget.NewEntry()
	g.components.watermarkEntry.SetPlaceHolder(locales[g.currentLocale].WatermarkPlaceholder)

	g.components.imageDirEntry = widget.NewEntry()
	g.components.imageDirEntry.SetPlaceHolder(locales[g.currentLocale].ImageDirPlaceholder)

	g.components.formatSelect = widget.NewSelect([]string{"jpg", "webp", "png"}, func(s string) {
		g.cfg.OutputFormat = s
	})
	g.components.formatSelect.SetSelected("jpg")

	g.components.qualityEntry = widget.NewEntry()
	g.components.qualityEntry.SetText(strconv.Itoa(g.cfg.Quality))
	g.components.qualityEntry.OnChanged = func(s string) {
		if q, err := strconv.Atoi(s); err == nil && q >= 1 && q <= 100 {
			g.cfg.Quality = q
		} else {
			g.components.qualityEntry.SetText(strconv.Itoa(g.cfg.Quality))
		}
	}

	g.components.widthEntry = widget.NewEntry()
	g.components.widthEntry.SetText(strconv.Itoa(g.cfg.MaxWidth))
	g.components.widthEntry.OnChanged = func(s string) {
		if w, err := strconv.Atoi(s); err == nil && w >= 1 {
			if g.components.watermarkEntry.Text != "" {
				if img, err := g.fileHandler.LoadImage(g.components.watermarkEntry.Text); err == nil && img != nil {
					if w > img.Bounds().Dx() {
						dialog.ShowInformation(locales[g.currentLocale].ErrorTitle, fmt.Sprintf(locales[g.currentLocale].WidthExceedsWatermark, img.Bounds().Dx()), g.window)
						g.components.widthEntry.SetText(strconv.Itoa(img.Bounds().Dx()))
						g.cfg.MaxWidth = img.Bounds().Dx()
						return
					}
				}
			}
			g.cfg.MaxWidth = w
		} else {
			g.components.widthEntry.SetText(strconv.Itoa(g.cfg.MaxWidth))
		}
	}

	g.components.heightEntry = widget.NewEntry()
	g.components.heightEntry.SetText(strconv.Itoa(g.cfg.MaxHeight))
	g.components.heightEntry.OnChanged = func(s string) {
		if h, err := strconv.Atoi(s); err == nil && h >= 1 {
			if g.components.watermarkEntry.Text != "" {
				if img, err := g.fileHandler.LoadImage(g.components.watermarkEntry.Text); err == nil && img != nil {
					if h > img.Bounds().Dy() {
						dialog.ShowInformation(locales[g.currentLocale].ErrorTitle, fmt.Sprintf(locales[g.currentLocale].HeightExceedsWatermark, img.Bounds().Dy()), g.window)
						g.components.heightEntry.SetText(strconv.Itoa(img.Bounds().Dy()))
						g.cfg.MaxHeight = img.Bounds().Dy()
						return
					}
				}
			}
			g.cfg.MaxHeight = h
		} else {
			g.components.heightEntry.SetText(strconv.Itoa(g.cfg.MaxHeight))
		}
	}

	g.components.targetSizeEntry = widget.NewEntry()
	g.components.targetSizeEntry.SetText("100")
	g.components.targetSizeEntry.OnChanged = func(s string) {
		if size, err := strconv.Atoi(s); err == nil && size >= 50 && size <= 5000 {
			fileio.NewImageProcessor().TargetSizeKB = size
		} else {
			g.components.targetSizeEntry.SetText("100")
		}
	}

	g.components.watermarkModeSelect = widget.NewSelect([]string{"crop", "resize"}, func(s string) {
		g.watermarkMode = s
	})
	g.components.watermarkModeSelect.SetSelected("crop")

	g.components.progress = widget.NewProgressBar()
	g.components.imageContainer = container.NewVBox()
	g.components.scrollContainer = container.NewVScroll(g.components.imageContainer)
	g.components.scrollContainer.SetMinSize(fyne.NewSize(400, 200))

	g.components.fileButton = g.createFileButton()
	g.components.folderButton = g.createFolderButton()

	g.components.languageSelect = widget.NewSelect([]string{"English", "Русский"}, func(s string) {
		if s == "English" {
			g.currentLocale = "en"
		} else {
			g.currentLocale = "ru"
		}
		g.updateLocale()
	})
	g.components.languageSelect.SetSelected("English")
}

func (g *GUI) updateLocale() {
	locale := locales[g.currentLocale]
	g.components.languageLabel.SetText(locale.LanguageLabel)
	g.components.widthLabel.SetText(locale.WidthLabel)
	g.components.heightLabel.SetText(locale.HeightLabel)
	g.components.watermarkLabel.SetText(locale.WatermarkLabel)
	g.components.imageDirLabel.SetText(locale.ImageDirLabel)
	g.components.formatLabel.SetText(locale.FormatLabel)
	g.components.qualityLabel.SetText(locale.QualityLabel)
	g.components.targetSizeLabel.SetText(locale.TargetSizeLabel)
	g.components.watermarkModeLabel.SetText(locale.WatermarkModeLabel)
	g.components.currentFileLabel.SetText(locale.CurrentFileLabel)
	g.components.webSizeHintLabel.SetText(locale.WebSizeHint)
	g.components.watermarkEntry.SetPlaceHolder(locale.WatermarkPlaceholder)
	g.components.imageDirEntry.SetPlaceHolder(locale.ImageDirPlaceholder)
	g.components.fileButton.SetText(locale.BrowseButton)
	g.components.folderButton.SetText(locale.BrowseFolderButton)
	g.window.Canvas().Refresh(g.window.Content())
}

func (g *GUI) createFileButton() *widget.Button {
	return widget.NewButton(locales[g.currentLocale].BrowseButton, func() {
		dialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				dialog.ShowInformation(locales[g.currentLocale].ErrorTitle, locales[g.currentLocale].FailedSelectWatermark, g.window)
				return
			}
			path := reader.URI().Path()
			if _, err := os.Stat(path); err != nil {
				dialog.ShowInformation(locales[g.currentLocale].ErrorTitle, fmt.Sprintf(locales[g.currentLocale].InvalidFile, err), g.window)
				return
			}
			g.components.watermarkEntry.SetText(path)
			if img, err := g.fileHandler.LoadImage(path); err == nil && img != nil {
				maxWidth := img.Bounds().Dx()
				maxHeight := img.Bounds().Dy()
				g.components.widthLabel.SetText(fmt.Sprintf("%s (100-%d):", locales[g.currentLocale].WidthLabel, maxWidth))
				g.components.heightLabel.SetText(fmt.Sprintf("%s (100-%d):", locales[g.currentLocale].HeightLabel, maxHeight))
				if g.cfg.MaxWidth > maxWidth {
					g.components.widthEntry.SetText(strconv.Itoa(maxWidth))
					g.cfg.MaxWidth = maxWidth
				}
				if g.cfg.MaxHeight > maxHeight {
					g.components.heightEntry.SetText(strconv.Itoa(maxHeight))
					g.cfg.MaxHeight = maxHeight
				}
			} else {
				g.components.widthLabel.SetText(locales[g.currentLocale].WidthLabel)
				g.components.heightLabel.SetText(locales[g.currentLocale].HeightLabel)
			}
		}, g.window)
		dialog.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg", ".webp"}))
		dialog.Show()
	})
}

func (g *GUI) createFolderButton() *widget.Button {
	return widget.NewButton(locales[g.currentLocale].BrowseFolderButton, func() {
		dialog := dialog.NewFolderOpen(func(reader fyne.ListableURI, err error) {
			if err != nil || reader == nil {
				dialog.ShowInformation(locales[g.currentLocale].ErrorTitle, locales[g.currentLocale].FailedSelectFolder, g.window)
				return
			}
			path := reader.Path()
			if _, err := os.Stat(path); err != nil {
				dialog.ShowInformation(locales[g.currentLocale].ErrorTitle, fmt.Sprintf(locales[g.currentLocale].InvalidFolder, err), g.window)
				return
			}
			g.components.imageDirEntry.SetText(path)
		}, g.window)
		dialog.Show()
	})
}

func (g *GUI) createProcessButton() *widget.Button {
	return widget.NewButton(locales[g.currentLocale].ProcessButton, func() {
		if g.components.watermarkEntry.Text == "" || g.components.imageDirEntry.Text == "" {
			dialog.ShowInformation(locales[g.currentLocale].ErrorTitle, locales[g.currentLocale].NoWatermarkOrFolder, g.window)
			return
		}

		if size, err := strconv.Atoi(g.components.targetSizeEntry.Text); err == nil && size >= 50 && size <= 5000 {
			fileio.NewImageProcessor().TargetSizeKB = size
		} else {
			g.components.targetSizeEntry.SetText("100")
			fileio.NewImageProcessor().TargetSizeKB = 100
		}

		g.components.progress.SetValue(0)
		g.components.imageContainer.Objects = nil
		g.components.currentFileLabel.SetText(locales[g.currentLocale].CurrentFileLabel)
		g.window.Canvas().Refresh(g.components.progress)
		g.window.Canvas().Refresh(g.components.currentFileLabel)

		processor, err := processor.NewImageProcessor(g.components.watermarkEntry.Text, g.cfg, g.fileHandler)
		if err != nil {
			dialog.ShowInformation(locales[g.currentLocale].ErrorTitle, fmt.Sprintf(locales[g.currentLocale].FailedInitProcessor, err), g.window)
			return
		}
		processor.WatermarkMode = g.watermarkMode
		err = processor.ProcessFolder(g.components.imageDirEntry.Text, g.cfg.OutputFormat, func(current, total int, img *canvas.Image, fileName string) {
			fyne.Do(func() {
				if total > 0 {
					g.components.progress.SetValue(float64(current) / float64(total))
					if img != nil {
						g.components.currentFileLabel.SetText(fmt.Sprintf("%s: %s", locales[g.currentLocale].CurrentFileLabel, fileName))
						fileInfo, err := os.Stat(filepath.Join("Images_watermarked", fileName+"."+g.cfg.OutputFormat))
						sizeKB := "N/A"
						if err == nil {
							sizeKB = fmt.Sprintf("%d KB", fileInfo.Size()/1024)
						}
						g.components.imageContainer.Add(container.NewVBox(
							widget.NewLabel(fmt.Sprintf("%s (%s)", fileName, sizeKB)),
							img,
						))
						g.components.scrollContainer.Refresh()
					}
					g.window.Canvas().Refresh(g.components.progress)
					g.window.Canvas().Refresh(g.components.currentFileLabel)
				}
			})
		})
		if err != nil {
			dialog.ShowInformation(locales[g.currentLocale].ErrorTitle, fmt.Sprintf(locales[g.currentLocale].ProcessingFailed, err), g.window)
		} else {
			g.components.currentFileLabel.SetText(locales[g.currentLocale].ProcessingDone)
			g.window.Canvas().Refresh(g.components.currentFileLabel)
		}
	})
}

func SetupGUI(window fyne.Window) {
	cfg := config.JpgConfig()
	fileHandler := &fileioHandler{}
	gui := NewGUI(window, cfg, fileHandler)
	window.SetIcon(Icon())
	gui.Setup()
}

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
