package main

import (
	"watermark/processor"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("GoIMGtool - Test")

	// folder watermark
	watermarkLabel := widget.NewLabel("Watermark file:")
	watermarkEntry := widget.NewEntry()
	watermarkEntry.SetPlaceHolder("Select watermark.png")

	// folder img
	imageDirLabel := widget.NewLabel("Image folder:")
	imageDirEntry := widget.NewEntry()
	imageDirEntry.SetPlaceHolder("Select image folder")

	// output
	formatLabel := widget.NewLabel("Output format:")
	formatSelect := widget.NewSelect([]string{"png", "webp"}, func(s string) {})
	formatSelect.SetSelected("png")

	progress := widget.NewProgressBar()

	// watermark button
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

	// folder button
	folderButton := widget.NewButton("Browse Folder...", func() {
		dialog := dialog.NewFolderOpen(func(reader fyne.ListableURI, err error) {
			if err != nil || reader == nil {
				return
			}
			imageDirEntry.SetText(reader.Path())
		}, w)
		dialog.Show()
	})

	// button
	processButton := widget.NewButton("Process", func() {
		if watermarkEntry.Text == "" || imageDirEntry.Text == "" {
			dialog.ShowInformation("Error", "Please select a watermark file and image folder!", w)
			return
		}
		progress.SetValue(0) // Progress drop
		processor, err := processor.NewImageProcessor(watermarkEntry.Text)
		if err != nil {
			dialog.ShowInformation("Error", err.Error(), w)
			return
		}
		err = processor.ProcessFolder(imageDirEntry.Text, formatSelect.Selected, func(current, total int) {
			if total > 0 {
				progress.SetValue(float64(current) / float64(total))
				w.Canvas().Refresh(progress) // Обновляем прогресс-бар
			}
		})
		if err != nil {
			dialog.ShowInformation("Error", err.Error(), w)
		}
	})

	// elem cont
	content := container.NewVBox(
		watermarkLabel,
		watermarkEntry,
		fileButton,
		imageDirLabel,
		imageDirEntry,
		folderButton,
		formatLabel,
		formatSelect,
		processButton,
		progress,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 350))
	w.ShowAndRun()
}
