package ui

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ShowAboutDialog displays the About dialog for the application.
func ShowAboutDialog(parent fyne.Window, dnglabVersion string) {
	appName := widget.NewLabelWithStyle(
		"DNGLab GUI",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	versionLabel := widget.NewLabel("Version 0.1.0")
	dnglabLabel := widget.NewLabel("dnglab: " + dnglabVersion)

	dnglabURL, _ := url.Parse("https://github.com/dnglab/dnglab")
	link := widget.NewHyperlink("dnglab on GitHub", dnglabURL)

	licenseLabel := widget.NewLabel("Copyright 2026 Juha Jantunen. MIT License.")

	content := container.NewVBox(
		appName,
		versionLabel,
		dnglabLabel,
		link,
		licenseLabel,
	)

	dialog.ShowCustom("About DNGLab GUI", "OK", content, parent)
}
