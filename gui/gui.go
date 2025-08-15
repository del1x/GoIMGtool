package gui

import (
	"github.com/del1x/GoIMGtool/processor"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func SetupGUI(w fyne.Window) {
	watermarkLabel := widget.NewLabel("Watermark file:")
	watermarkEntry := widget.NewEntry()
	watermarkEntry.SetPlaceHolder("Select watermark.png")

	imageDirLabel := widget.NewLabel("Image folder:")
	imageDirEntry := widget.NewEntry()
	imageDirEntry.SetPlaceHolder("Select image folder")

	formatLabel := widget.NewLabel("Output format:")
	formatSelect := widget.NewSelect([]string{"png", "webp"}, func(s string) {})
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

	processButton := widget.NewButton("Process", func() {
		if watermarkEntry.Text == "" || imageDirEntry.Text == "" {
			dialog.ShowInformation("Error", "Please select a watermark file and image folder!", w)
			return
		}
		progress.SetValue(0)
		processor, err := processor.NewImageProcessor(watermarkEntry.Text)
		if err != nil {
			dialog.ShowInformation("Error", err.Error(), w)
			return
		}
		err = processor.ProcessFolder(imageDirEntry.Text, formatSelect.Selected, func(current, total int) {
			if total > 0 {
				progress.SetValue(float64(current) / float64(total))
				w.Canvas().Refresh(progress)
			}
		})
		if err != nil {
			dialog.ShowInformation("Error", err.Error(), w)
		}
	})

	content := container.NewVBox(
		watermarkLabel, watermarkEntry, fileButton,
		imageDirLabel, imageDirEntry, folderButton,
		formatLabel, formatSelect,
		processButton, progress,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 350))
}
