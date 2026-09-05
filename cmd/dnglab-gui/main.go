package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("DNGLab GUI")

	w.SetContent(container.NewVBox(
		widget.NewLabel("DNGLab GUI"),
	))

	w.ShowAndRun()
}
