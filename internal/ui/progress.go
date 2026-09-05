package ui

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ProgressDialog shows conversion progress in a separate window.
type ProgressDialog struct {
	win       fyne.Window
	parent    fyne.Window
	cancelFn  context.CancelFunc
	cancelled bool

	fileLabel   *widget.Label
	progressBar *widget.ProgressBarInfinite
	logEntry    *widget.Entry
}

// NewProgressDialog creates a new progress dialog. Call Show() to display it.
func NewProgressDialog(a fyne.App, parent fyne.Window) *ProgressDialog {
	pd := &ProgressDialog{
		parent: parent,
	}

	pd.win = a.NewWindow("Converting...")
	pd.win.Resize(fyne.NewSize(500, 350))

	pd.fileLabel = widget.NewLabel("Starting...")

	pd.progressBar = widget.NewProgressBarInfinite()

	pd.logEntry = widget.NewMultiLineEntry()
	pd.logEntry.SetPlaceHolder("")
	pd.logEntry.Disable()
	pd.logEntry.Wrapping = fyne.TextWrapOff

	cancelBtn := widget.NewButton("Cancel", func() {
		pd.cancelled = true
		if pd.cancelFn != nil {
			pd.cancelFn()
		}
	})

	content := container.NewBorder(
		container.NewVBox(
			pd.fileLabel,
			pd.progressBar,
		),
		cancelBtn,
		nil, nil,
		container.NewVScroll(pd.logEntry),
	)

	pd.win.SetContent(content)

	return pd
}

// SetCancelFunc registers the context cancel function for the Cancel button.
func (pd *ProgressDialog) SetCancelFunc(cancel context.CancelFunc) {
	pd.cancelFn = cancel
}

// Show displays the progress window.
func (pd *ProgressDialog) Show() {
	pd.progressBar.Start()
	pd.win.Show()
}

// Update refreshes the progress display. Safe to call from any goroutine.
//
// Note: fyne.Do (for explicit main-thread dispatch) is available only from
// Fyne v2.6.0. This project targets Fyne v2.4.0 where widget setter methods
// (SetText) use internal property locks and are safe to call from a goroutine.
// When upgrading to Fyne v2.6+, wrap these calls in fyne.Do.
func (pd *ProgressDialog) Update(file string, current, total int, logLine string) {
	pd.fileLabel.SetText(fmt.Sprintf("Converting %d files...", total))

	if logLine != "" {
		existing := pd.logEntry.Text
		if existing != "" {
			existing += "\n"
		}
		pd.logEntry.SetText(existing + logLine)
	}
}

// Complete closes the progress window and shows a result dialog on the parent window.
//
// Note: see Update for the Fyne v2.4.0 goroutine-safety note. This method also
// calls pd.win.Hide(), dialog.ShowInformation, and dialog.ShowError from the
// same goroutine. Whether Fyne v2.4.0 marshals window hide and dialog
// construction to the main thread internally is not explicitly guaranteed by
// the Fyne v2.4 API. When upgrading to Fyne v2.6+, wrap the entire body in
// fyne.Do to ensure all Fyne state changes happen on the main thread.
func (pd *ProgressDialog) Complete(err error) {
	pd.progressBar.Stop()
	pd.win.Hide()
	if pd.cancelled {
		dialog.ShowInformation("Conversion cancelled", "The conversion was cancelled.", pd.parent)
		return
	}
	if err != nil {
		dialog.ShowError(err, pd.parent)
		return
	}
	dialog.ShowInformation("Conversion complete", "All files have been converted successfully.", pd.parent)
}

// IsCancelled reports whether the user clicked Cancel.
func (pd *ProgressDialog) IsCancelled() bool {
	return pd.cancelled
}
