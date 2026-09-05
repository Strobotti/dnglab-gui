package converter

import (
	"testing"
)

func TestIsSupportedRawFile(t *testing.T) {
	tests := []struct {
		filename string
		want     bool
	}{
		// Uppercase extensions
		{"photo.CR3", true},
		{"photo.NEF", true},
		{"photo.ARW", true},

		// Lowercase extensions
		{"photo.cr3", true},
		{"photo.nef", true},
		{"photo.arw", true},

		// Mixed case
		{"photo.Cr3", true},

		// Non-RAW extensions
		{"photo.jpg", false},
		{"photo.png", false},
		{"photo.dng", false}, // DNG is the output format, not a valid input
		{"photo.tiff", false},
		{"photo.jpeg", false},

		// No extension
		{"photoonly", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got := IsSupportedRawFile(tt.filename)
			if got != tt.want {
				t.Errorf("IsSupportedRawFile(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}
