package ui

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/strobotti/dnglab-gui/internal/converter"
)

const dnglabReleasesURL = "https://github.com/dnglab/dnglab/releases"

// CheckDNGLabOnStart looks for the dnglab binary once the application has
// started. If it is missing, the user is asked whether to open the dnglab
// download page (the application then closes) or cancel (the application closes).
func CheckDNGLabOnStart(a fyne.App, w fyne.Window) {
	a.Lifecycle().SetOnStarted(func() {
		dl := converter.DNGLab{}
		if err := dl.DetectBinary(); err == nil {
			return
		}
		showMissingDNGLab(a, w)
	})
}

func showMissingDNGLab(a fyne.App, w fyne.Window) {
	msg := widget.NewLabel(
		"The dnglab binary was not found next to this application or on your PATH.\n\n" +
			"Open the dnglab releases page to download it? The application will close.")
	msg.Wrapping = fyne.TextWrapWord

	d := dialog.NewCustomConfirm(
		"dnglab not found",
		"Open download page",
		"Cancel",
		msg,
		func(open bool) {
			if open {
				if u, err := url.Parse(dnglabReleasesURL); err == nil {
					_ = a.OpenURL(u)
				}
			}
			a.Quit()
		},
		w,
	)
	d.Show()
}
