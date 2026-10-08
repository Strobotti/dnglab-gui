package ui

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/strobotti/dnglab-gui/internal/converter"
)

// ProgressDialog shows conversion progress in a separate window, similar to
// the Adobe DNG Converter: processed/total counts, success and error counts,
// the current file, elapsed time and an estimate of the remaining time.
type ProgressDialog struct {
	win    fyne.Window
	parent fyne.Window
	total  int

	mu        sync.Mutex
	cancelFn  context.CancelFunc
	cancelled bool
	processed int
	converted int
	skipped   int
	failed    int
	current   string
	start     time.Time
	stopTick  chan struct{}

	processedLabel *widget.Label
	progressBar    *widget.ProgressBar
	convertedLabel *widget.Label
	skippedLabel   *widget.Label
	failedLabel    *widget.Label
	currentLabel   *widget.Label
	timeLabel      *widget.Label
	logEntry       *widget.Entry
	cancelBtn      *widget.Button
}

// NewProgressDialog creates a new progress dialog for total files. Call Show() to display it.
func NewProgressDialog(a fyne.App, parent fyne.Window, total int) *ProgressDialog {
	pd := &ProgressDialog{
		parent: parent,
		total:  total,
	}

	pd.win = a.NewWindow("Converting...")
	pd.win.Resize(fyne.NewSize(560, 430))

	pd.processedLabel = widget.NewLabelWithStyle(
		fmt.Sprintf("Processed: 0 of %d", total), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	pd.progressBar = widget.NewProgressBar()
	pd.convertedLabel = widget.NewLabel("Converted: 0")
	pd.skippedLabel = widget.NewLabel("Skipped (already exist): 0")
	pd.failedLabel = widget.NewLabel("Errors: 0")
	pd.currentLabel = widget.NewLabel("Current file: starting...")
	pd.currentLabel.Truncation = fyne.TextTruncateEllipsis
	pd.timeLabel = widget.NewLabel("Elapsed: 00:00   Remaining: estimating...")

	pd.logEntry = widget.NewMultiLineEntry()
	pd.logEntry.Disable()
	pd.logEntry.Wrapping = fyne.TextWrapOff

	pd.cancelBtn = widget.NewButton("Cancel", func() {
		pd.mu.Lock()
		pd.cancelled = true
		cancel := pd.cancelFn
		pd.mu.Unlock()
		if cancel != nil {
			cancel()
		}
	})

	header := container.NewVBox(
		pd.processedLabel,
		pd.progressBar,
		container.NewGridWithColumns(3, pd.convertedLabel, pd.skippedLabel, pd.failedLabel),
		pd.currentLabel,
		pd.timeLabel,
		widget.NewSeparator(),
	)

	content := container.NewBorder(
		header,
		container.NewHBox(layout.NewSpacer(), pd.cancelBtn),
		nil, nil,
		container.NewVScroll(pd.logEntry),
	)

	pd.win.SetContent(content)
	return pd
}

// SetCancelFunc registers the context cancel function for the Cancel button.
func (pd *ProgressDialog) SetCancelFunc(cancel context.CancelFunc) {
	pd.mu.Lock()
	defer pd.mu.Unlock()
	pd.cancelFn = cancel
}

// Show displays the progress window and starts the elapsed-time ticker.
func (pd *ProgressDialog) Show() {
	pd.mu.Lock()
	pd.start = time.Now()
	pd.stopTick = make(chan struct{})
	pd.mu.Unlock()

	pd.win.Show()
	pd.refresh()

	go func(stop <-chan struct{}) {
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for {
			select {
			case <-stop:
				return
			case <-t.C:
				pd.refresh()
			}
		}
	}(pd.stopTick)
}

// HandleEvent processes one line of dnglab output. Safe to call from any goroutine.
func (pd *ProgressDialog) HandleEvent(ev converter.Event) {
	pd.mu.Lock()
	switch {
	case ev.Converted:
		pd.processed++
		pd.converted++
		pd.current = filepath.Base(ev.Source)
	case ev.Skipped:
		pd.processed++
		pd.skipped++
		pd.current = filepath.Base(ev.Source)
	case ev.Failed:
		pd.processed++
		pd.failed++
		pd.current = filepath.Base(ev.Source)
	}
	pd.mu.Unlock()

	if ev.Line != "" {
		pd.appendLog(ev.Line)
	}
	pd.refresh()
}

// Complete stops the ticker, closes the progress window and shows the result
// on the parent window. err is set only for failures that stopped the run.
func (pd *ProgressDialog) Complete(res converter.Result, err error) {
	pd.mu.Lock()
	if pd.stopTick != nil {
		close(pd.stopTick)
		pd.stopTick = nil
	}
	cancelled := pd.cancelled
	pd.mu.Unlock()

	fyne.Do(func() {
		pd.win.Hide()
		processed := res.Converted + res.Skipped + len(res.Failed)
		switch {
		case cancelled:
			dialog.ShowInformation("Conversion cancelled",
				fmt.Sprintf("Cancelled after %d of %d files (%d converted, %d skipped, %d errors).",
					processed, pd.total, res.Converted, res.Skipped, len(res.Failed)), pd.parent)
		case err != nil:
			dialog.ShowError(err, pd.parent)
		case len(res.Failed) == 0 && res.Skipped == 0:
			dialog.ShowInformation("Conversion complete",
				fmt.Sprintf("Converted %d of %d files.", res.Converted, pd.total), pd.parent)
		default:
			showResultDialog(pd.parent, pd.total, res)
		}
	})
}

// showResultDialog shows the summary of a run that had skipped or failed files.
// Failed files are listed in a scrollable list, so large runs stay readable.
func showResultDialog(parent fyne.Window, total int, res converter.Result) {
	summary := widget.NewLabelWithStyle(
		fmt.Sprintf("Converted: %d    Skipped: %d    Failed: %d    Total: %d",
			res.Converted, res.Skipped, len(res.Failed), total),
		fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	top := container.NewVBox(summary)
	if res.Skipped > 0 {
		top.Add(widget.NewLabel("Skipped files already exist in the destination and were not overwritten."))
	}

	var body fyne.CanvasObject
	if len(res.Failed) == 0 {
		body = widget.NewLabel("No files failed.")
	} else {
		top.Add(widget.NewLabelWithStyle("Failed files:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		list := widget.NewList(
			func() int { return len(res.Failed) },
			func() fyne.CanvasObject {
				l := widget.NewLabel("")
				l.Truncation = fyne.TextTruncateEllipsis
				return l
			},
			func(id widget.ListItemID, obj fyne.CanvasObject) {
				f := res.Failed[id]
				obj.(*widget.Label).SetText(filepath.Base(f.Path) + ": " + strings.TrimPrefix(f.Reason, "Status: Failed: "))
			},
		)
		body = list
	}

	content := container.NewBorder(top, nil, nil, nil, body)
	d := dialog.NewCustom("Conversion finished", "OK", content, parent)
	d.Resize(fyne.NewSize(720, 480))
	d.Show()
}

// IsCancelled reports whether the user clicked Cancel.
func (pd *ProgressDialog) IsCancelled() bool {
	pd.mu.Lock()
	defer pd.mu.Unlock()
	return pd.cancelled
}

// refresh redraws all counters from the current state.
func (pd *ProgressDialog) refresh() {
	pd.mu.Lock()
	processed, converted, skipped, failed := pd.processed, pd.converted, pd.skipped, pd.failed
	current := pd.current
	elapsed := time.Duration(0)
	if !pd.start.IsZero() {
		elapsed = time.Since(pd.start)
	}
	pd.mu.Unlock()

	remaining := "estimating..."
	if processed > 0 && processed < pd.total {
		eta := time.Duration(float64(elapsed) / float64(processed) * float64(pd.total-processed))
		remaining = formatDuration(eta)
	} else if processed >= pd.total && pd.total > 0 {
		remaining = "0:00"
	}

	fraction := 0.0
	if pd.total > 0 {
		fraction = float64(processed) / float64(pd.total)
	}
	if current == "" {
		current = "starting..."
	}

	fyne.Do(func() {
		pd.processedLabel.SetText(fmt.Sprintf("Processed: %d of %d", processed, pd.total))
		pd.progressBar.SetValue(fraction)
		pd.convertedLabel.SetText(fmt.Sprintf("Converted: %d", converted))
		pd.skippedLabel.SetText(fmt.Sprintf("Skipped (already exist): %d", skipped))
		pd.failedLabel.SetText(fmt.Sprintf("Errors: %d", failed))
		pd.currentLabel.SetText("Current file: " + current)
		pd.timeLabel.SetText(fmt.Sprintf("Elapsed: %s   Remaining: %s", formatDuration(elapsed), remaining))
	})
}

func (pd *ProgressDialog) appendLog(line string) {
	fyne.Do(func() {
		existing := pd.logEntry.Text
		if existing != "" {
			existing += "\n"
		}
		pd.logEntry.SetText(existing + line)
	})
}

// formatDuration renders d as m:ss, or h:mm:ss when it is an hour or longer.
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}
