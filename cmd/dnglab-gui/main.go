package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/strobotti/dnglab-gui/internal/ui"
)

func main() {
	a := app.NewWithID("com.github.strobotti.dnglab-gui")
	w := ui.NewMainWindow(a)
	w.ShowAndRun()
}
