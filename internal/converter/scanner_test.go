package converter

import (
	"os"
	"path/filepath"
	"sort"
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

// TestScanForRawFiles exercises ScanForRawFiles with a real temp-directory
// fixture, including recursive and non-recursive modes and non-RAW file
// filtering.
func TestScanForRawFiles(t *testing.T) {
	// Build a fixture tree:
	//   <root>/
	//     a.CR3
	//     b.jpg          (ignored)
	//     sub/
	//       c.NEF
	//       d.txt        (ignored)

	root := t.TempDir()
	sub := filepath.Join(root, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatalf("creating sub dir: %v", err)
	}

	touch := func(path string) {
		t.Helper()
		if err := os.WriteFile(path, nil, 0o644); err != nil {
			t.Fatalf("creating fixture file %s: %v", path, err)
		}
	}

	fileA := filepath.Join(root, "a.CR3")
	fileC := filepath.Join(sub, "c.NEF")
	touch(fileA)
	touch(filepath.Join(root, "b.jpg"))
	touch(fileC)
	touch(filepath.Join(sub, "d.txt"))

	tests := []struct {
		name      string
		recursive bool
		wantFiles []string // absolute paths, sorted
	}{
		{
			name:      "non-recursive only finds top-level RAW files",
			recursive: false,
			wantFiles: []string{fileA},
		},
		{
			name:      "recursive finds RAW files in all subdirectories",
			recursive: true,
			wantFiles: []string{fileA, fileC},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ScanForRawFiles(root, tt.recursive)
			if err != nil {
				t.Fatalf("ScanForRawFiles returned unexpected error: %v", err)
			}

			// Sort both slices for a stable comparison.
			sort.Strings(got)
			want := make([]string, len(tt.wantFiles))
			copy(want, tt.wantFiles)
			sort.Strings(want)

			if len(got) != len(want) {
				t.Fatalf("got %d files %v, want %d files %v", len(got), got, len(want), want)
			}
			for i := range want {
				if got[i] != want[i] {
					t.Errorf("file[%d]: got %q, want %q", i, got[i], want[i])
				}
			}
		})
	}
}

