package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/strobotti/dnglab-gui/internal/ui"
)

func main() {
	a := app.New()
	w := ui.NewMainWindow(a)
	w.ShowAndRun()
}
