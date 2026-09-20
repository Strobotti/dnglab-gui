package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/strobotti/dnglab-gui/internal/assets"
	"github.com/strobotti/dnglab-gui/internal/ui"
)

func main() {
	a := app.NewWithID("com.github.strobotti.dnglab-gui")
	a.SetIcon(assets.AppIcon())
	w := ui.NewMainWindow(a)
	w.ShowAndRun()
}
