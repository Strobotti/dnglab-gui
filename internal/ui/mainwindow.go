package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/strobotti/dnglab-gui/internal/converter"
)

// NewMainWindow creates and returns the application's main window.
func NewMainWindow(a fyne.App) fyne.Window {
	w := a.NewWindow("DNGLab GUI")
	w.Resize(fyne.NewSize(640, 560))

	opts := converter.DefaultConvertOptions()

	// ── Section 1: Select images to convert ─────────────────────────────────

	sec1Label := widget.NewLabelWithStyle(
		"1  Select images to convert",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	inputPathLabel := widget.NewLabel("No images have been selected")

	// Convert button declared early so the folder selector can enable it.
	convertBtn := widget.NewButton("Convert", func() {
		dialog.ShowInformation("Convert", "Conversion will be implemented in the next step", w)
	})
	convertBtn.Disable()

	selectInputBtn := widget.NewButton("Select Folder...", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			opts.InputPath = uri.Path()
			inputPathLabel.SetText(uri.Path())
			convertBtn.Enable()
		}, w)
	})

	recursiveCheck := widget.NewCheck("Include images in subfolders", func(checked bool) {
		opts.Recursive = checked
	})

	skipCheck := widget.NewCheck("Skip source image if destination image already exists", func(checked bool) {
		opts.SkipExisting = checked
	})
	skipCheck.SetChecked(true) // matches DefaultConvertOptions

	section1 := container.NewVBox(
		sec1Label,
		container.NewHBox(selectInputBtn, inputPathLabel),
		recursiveCheck,
		skipCheck,
		widget.NewSeparator(),
	)

	// ── Section 2: Select location to save converted images ──────────────────

	sec2Label := widget.NewLabelWithStyle(
		"2  Select location to save converted images",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	outputPathLabel := widget.NewLabel("")

	outputFolderBtn := widget.NewButton("Select Folder...", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			opts.OutputPath = uri.Path()
			outputPathLabel.SetText(uri.Path())
		}, w)
	})
	outputFolderBtn.Disable()

	locationSelect := widget.NewSelect(
		[]string{"Save in Same Location", "Select Folder..."},
		func(selected string) {
			switch selected {
			case "Save in Same Location":
				opts.OutputPath = ""
				outputPathLabel.SetText("")
				outputFolderBtn.Disable()
			case "Select Folder...":
				outputFolderBtn.Enable()
			}
		},
	)
	locationSelect.SetSelected("Save in Same Location")

	section2 := container.NewVBox(
		sec2Label,
		locationSelect,
		container.NewHBox(outputFolderBtn, outputPathLabel),
		widget.NewSeparator(),
	)

	// ── Section 3: Conversion Options ───────────────────────────────────────

	sec3Label := widget.NewLabelWithStyle(
		"3  Conversion Options",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	compressionRadio := widget.NewRadioGroup([]string{"Lossless", "Uncompressed"}, func(s string) {
		opts.Compression = strings.ToLower(s)
	})
	compressionRadio.SetSelected("Lossless")

	cropMap := map[string]string{
		"Best":        "best",
		"Active Area": "activearea",
		"None":        "none",
	}
	cropRadio := widget.NewRadioGroup([]string{"Best", "Active Area", "None"}, func(s string) {
		if v, ok := cropMap[s]; ok {
			opts.Crop = v
		}
	})
	cropRadio.SetSelected("Best")

	embedRawCheck := widget.NewCheck("Embed original RAW file", func(checked bool) {
		opts.EmbedRaw = checked
	})
	embedRawCheck.SetChecked(true)

	previewCheck := widget.NewCheck("Include preview image", func(checked bool) {
		opts.Preview = checked
	})
	previewCheck.SetChecked(true)

	thumbnailCheck := widget.NewCheck("Include thumbnail", func(checked bool) {
		opts.Thumbnail = checked
	})
	thumbnailCheck.SetChecked(true)

	section3 := container.NewVBox(
		sec3Label,
		widget.NewLabel("Compression:"),
		compressionRadio,
		widget.NewLabel("Crop:"),
		cropRadio,
		embedRawCheck,
		previewCheck,
		thumbnailCheck,
		widget.NewSeparator(),
	)

	// ── Section 4: Metadata ──────────────────────────────────────────────────

	sec4Label := widget.NewLabelWithStyle(
		"4  Metadata",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	artistEntry := widget.NewEntry()
	artistEntry.OnChanged = func(s string) {
		opts.Artist = s
	}

	section4 := container.NewVBox(
		sec4Label,
		container.NewHBox(widget.NewLabel("Artist:"), artistEntry),
		widget.NewSeparator(),
	)

	// ── Bottom button bar ────────────────────────────────────────────────────

	aboutBtn := widget.NewButton("About", func() {
		dialog.ShowInformation("About", "DNGLab GUI v0.1.0", w)
	})

	quitBtn := widget.NewButton("Quit", func() {
		a.Quit()
	})

	buttonBar := container.NewHBox(
		aboutBtn,
		quitBtn,
		layout.NewSpacer(),
		convertBtn,
	)

	// ── Assemble the full layout ─────────────────────────────────────────────

	content := container.NewVScroll(
		container.NewVBox(
			section1,
			section2,
			section3,
			section4,
			buttonBar,
		),
	)

	w.SetContent(content)

	return w
}
