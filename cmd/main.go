package main

import (
	"github.com/del1x/GoIMGtool/gui"

	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("GoIMGtool")
	gui.SetupGUI(w)
	w.ShowAndRun()
}
