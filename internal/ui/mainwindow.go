package ui

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/strobotti/dnglab-gui/internal/config"
	"github.com/strobotti/dnglab-gui/internal/converter"
)

// browseFolder opens the desktop's native folder picker so it can be resized
// and moved, and shows removable media. On Linux this uses zenity (GTK file
// chooser). If zenity is not installed it falls back to the Fyne folder dialog.
// The native picker blocks, so it runs in a goroutine; onPick runs on the UI thread.
func browseFolder(w fyne.Window, title, startDir string, onPick func(path string)) {
	if _, err := exec.LookPath("zenity"); err != nil {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			onPick(uri.Path())
		}, w)
		return
	}

	go func() {
		args := []string{"--file-selection", "--directory", "--title=" + title}
		if startDir != "" {
			args = append(args, "--filename="+strings.TrimRight(startDir, "/")+"/")
		}
		out, err := exec.Command("zenity", args...).Output()
		if err != nil {
			// zenity exits non-zero when the user cancels.
			return
		}
		path := strings.TrimSpace(string(out))
		if path == "" {
			return
		}
		fyne.Do(func() { onPick(path) })
	}()
}

// NewMainWindow creates and returns the application's main window.
func NewMainWindow(a fyne.App) fyne.Window {
	w := a.NewWindow("DNGLab GUI")
	w.Resize(fyne.NewSize(640, 780))

	// Load persisted settings.
	settings := config.Load(a.Preferences())

	opts := converter.DefaultConvertOptions()
	// Apply saved settings to opts.
	opts.Compression = settings.Compression
	opts.Crop = settings.Crop
	opts.EmbedRaw = settings.EmbedRaw
	opts.Preview = settings.Preview
	opts.Thumbnail = settings.Thumbnail
	opts.SkipExisting = settings.SkipExisting
	opts.Recursive = settings.Recursive
	opts.Artist = settings.Artist
	opts.InputPath = settings.LastInputDir
	opts.OutputPath = settings.LastOutputDir

	// ── Section 1: Select images to convert ─────────────────────────────────

	sec1Label := widget.NewLabelWithStyle(
		"1. Select images to convert",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	inputPathText := settings.LastInputDir
	if inputPathText == "" {
		inputPathText = "No images have been selected"
	}
	inputPathLabel := widget.NewLabel(inputPathText)

	// Convert button declared early so the folder selector can enable it.
	// Convert is enabled only when an input folder is selected and the settings
	// differ from those of the last successful conversion.
	var lastRunKey string
	var convertBtn *widget.Button
	updateConvertButton := func() {
		if convertBtn == nil {
			return
		}
		key := fmt.Sprintf("%+v", opts)
		if opts.InputPath != "" && (lastRunKey == "" || key != lastRunKey) {
			convertBtn.Enable()
		} else {
			convertBtn.Disable()
		}
	}
	convertBtn = widget.NewButton("Convert", func() {
		// 1. Detect dnglab binary.
		dl := converter.DNGLab{}
		if err := dl.DetectBinary(); err != nil {
			dialog.ShowError(fmt.Errorf("dnglab not found: %v\n\nPlease install dnglab and ensure it is on your PATH", err), w)
			return
		}

		// 2. Scan for RAW files.
		files, err := converter.ScanForRawFiles(opts.InputPath, opts.Recursive)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if len(files) == 0 {
			dialog.ShowInformation("No files found", "No supported RAW files were found in the selected folder.", w)
			return
		}

		// 3. Show progress dialog.
		ctx, cancel := context.WithCancel(context.Background())
		pd := NewProgressDialog(a, w, len(files))
		pd.SetCancelFunc(cancel)
		pd.Show()

		// 4. Disable Convert button during conversion.
		convertBtn.Disable()

		// 5. Save settings before launching so an abrupt exit does not lose them.
		config.Save(a.Preferences(), settings)

		// 6. Run conversion in a goroutine; dnglab's output drives the progress dialog.
		runKey := fmt.Sprintf("%+v", opts)
		go func() {
			res, convErr := dl.Convert(ctx, opts, len(files), pd.HandleEvent)
			cancel()
			pd.Complete(res, convErr)
			fyne.Do(func() {
				if convErr == nil && len(res.Failed) == 0 && !pd.IsCancelled() {
					lastRunKey = runKey
				}
				updateConvertButton()
			})
		}()
	})
	updateConvertButton()

	selectInputBtn := widget.NewButton("Select Folder...", func() {
		browseFolder(w, "Select source folder", settings.LastInputDir, func(path string) {
			opts.InputPath = path
			settings.LastInputDir = path
			inputPathLabel.SetText(path)
			updateConvertButton()
		})
	})

	recursiveCheck := widget.NewCheck("Include images in subfolders", func(checked bool) {
		opts.Recursive = checked
		settings.Recursive = checked
		updateConvertButton()
	})
	recursiveCheck.SetChecked(settings.Recursive)

	skipCheck := widget.NewCheck("Skip source image if destination image already exists", func(checked bool) {
		opts.SkipExisting = checked
		updateConvertButton()
	})
	skipCheck.SetChecked(opts.SkipExisting)

	section1 := container.NewVBox(
		sec1Label,
		container.NewHBox(selectInputBtn, inputPathLabel),
		recursiveCheck,
		skipCheck,
		widget.NewSeparator(),
	)

	// ── Section 2: Select location to save converted images ──────────────────

	sec2Label := widget.NewLabelWithStyle(
		"2. Select location to save converted images",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	outputPathLabel := widget.NewLabel(opts.OutputPath)

	outputFolderBtn := widget.NewButton("Select Folder...", func() {
		browseFolder(w, "Select output folder", settings.LastOutputDir, func(path string) {
			opts.OutputPath = path
			settings.LastOutputDir = path
			outputPathLabel.SetText(path)
			updateConvertButton()
		})
	})

	// Determine initial state of the output folder button.
	initialLocationSelection := "Save in Same Location"
	if settings.OutputMode == "folder" {
		initialLocationSelection = "Select Folder..."
		outputFolderBtn.Enable()
	} else {
		outputFolderBtn.Disable()
	}

	locationSelect := widget.NewSelect(
		[]string{"Save in Same Location", "Select Folder..."},
		func(selected string) {
			switch selected {
			case "Save in Same Location":
				opts.OutputPath = ""
				outputPathLabel.SetText("")
				outputFolderBtn.Disable()
				settings.OutputMode = "same"
			case "Select Folder...":
				outputFolderBtn.Enable()
				settings.OutputMode = "folder"
			}
			updateConvertButton()
		},
	)
	locationSelect.SetSelected(initialLocationSelection)

	section2 := container.NewVBox(
		sec2Label,
		container.NewBorder(nil, nil, nil, outputFolderBtn, locationSelect),
		outputPathLabel,
		widget.NewSeparator(),
	)

	// ── Section 3: Conversion Options ───────────────────────────────────────

	sec3Label := widget.NewLabelWithStyle(
		"3. Conversion Options",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	compressionRadio := widget.NewRadioGroup([]string{"Lossless", "Uncompressed"}, func(s string) {
		opts.Compression = strings.ToLower(s)
		settings.Compression = opts.Compression
		updateConvertButton()
	})
	// Restore saved compression selection.
	switch opts.Compression {
	case "uncompressed":
		compressionRadio.SetSelected("Uncompressed")
	default:
		compressionRadio.SetSelected("Lossless")
	}

	// cropMap translates UI display strings to dnglab --crop flag values.
	// The dnglab CLI accepts "best", "activearea" (no hyphen), and "none" --
	// this has been verified against the official dnglab documentation at
	// https://github.com/dnglab/dnglab (possible values: best, activearea, none).
	cropMap := map[string]string{
		"Best":        "best",
		"Active Area": "activearea",
		"None":        "none",
	}
	cropRadio := widget.NewRadioGroup([]string{"Best", "Active Area", "None"}, func(s string) {
		if v, ok := cropMap[s]; ok {
			opts.Crop = v
			settings.Crop = v
			updateConvertButton()
		}
	})
	// Restore saved crop selection.
	switch opts.Crop {
	case "activearea":
		cropRadio.SetSelected("Active Area")
	case "none":
		cropRadio.SetSelected("None")
	default:
		cropRadio.SetSelected("Best")
	}

	embedRawCheck := widget.NewCheck("Embed original RAW file", func(checked bool) {
		opts.EmbedRaw = checked
		settings.EmbedRaw = checked
		updateConvertButton()
	})
	embedRawCheck.SetChecked(opts.EmbedRaw)

	previewCheck := widget.NewCheck("Include preview image", func(checked bool) {
		opts.Preview = checked
		settings.Preview = checked
		updateConvertButton()
	})
	previewCheck.SetChecked(opts.Preview)

	thumbnailCheck := widget.NewCheck("Include thumbnail", func(checked bool) {
		opts.Thumbnail = checked
		settings.Thumbnail = checked
		updateConvertButton()
	})
	thumbnailCheck.SetChecked(opts.Thumbnail)

	section3 := container.NewVBox(
		sec3Label,
		container.NewHBox(widget.NewLabel("Compression:"), compressionRadio),
		container.NewHBox(widget.NewLabel("Crop:"), cropRadio),
		embedRawCheck,
		previewCheck,
		thumbnailCheck,
		widget.NewSeparator(),
	)

	// ── Section 4: Metadata ──────────────────────────────────────────────────

	sec4Label := widget.NewLabelWithStyle(
		"4. Metadata",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	artistEntry := widget.NewEntry()
	artistEntry.SetText(opts.Artist)
	artistEntry.OnChanged = func(s string) {
		opts.Artist = s
		settings.Artist = s
		updateConvertButton()
	}

	section4 := container.NewVBox(
		sec4Label,
		container.NewBorder(nil, nil, widget.NewLabel("Artist:"), nil, artistEntry),
		widget.NewSeparator(),
	)

	// ── Bottom button bar ────────────────────────────────────────────────────

	aboutBtn := widget.NewButton("About", func() {
		dl := converter.DNGLab{}
		_ = dl.DetectBinary()
		ver, err := dl.Version()
		if err != nil {
			ver = "not detected"
		}
		ShowAboutDialog(w, ver)
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

	// Save settings when the window closes.
	w.SetOnClosed(func() {
		config.Save(a.Preferences(), settings)
	})

	return w
}
