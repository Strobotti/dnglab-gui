package converter

import (
	"testing"
)

func contains(args []string, value string) bool {
	for _, a := range args {
		if a == value {
			return true
		}
	}
	return false
}

func containsSeq(args []string, a, b string) bool {
	for i := 0; i < len(args)-1; i++ {
		if args[i] == a && args[i+1] == b {
			return true
		}
	}
	return false
}

// TestToArgs_Defaults verifies that default options produce correct flags and
// omit -f, -r, and --artist.
func TestToArgs_Defaults(t *testing.T) {
	opts := DefaultConvertOptions()
	args := opts.ToArgs()

	if contains(args, "-f") {
		t.Error("default options should NOT contain -f (SkipExisting=true means no override)")
	}
	if contains(args, "-r") {
		t.Error("default options should NOT contain -r (Recursive=false)")
	}
	if contains(args, "--artist") {
		t.Error("default options should NOT contain --artist (Artist is empty)")
	}
	if !containsSeq(args, "-c", "lossless") {
		t.Error("default options should contain '-c lossless'")
	}
	if !containsSeq(args, "--crop", "best") {
		t.Error("default options should contain '--crop best'")
	}
	if !containsSeq(args, "--embed-raw", "true") {
		t.Error("default options should contain '--embed-raw true'")
	}
	if !containsSeq(args, "--dng-preview", "true") {
		t.Error("default options should contain '--dng-preview true'")
	}
	if !containsSeq(args, "--dng-thumbnail", "true") {
		t.Error("default options should contain '--dng-thumbnail true'")
	}
}

// TestToArgs_Recursive verifies that Recursive=true produces -r.
func TestToArgs_Recursive(t *testing.T) {
	opts := DefaultConvertOptions()
	opts.Recursive = true
	args := opts.ToArgs()

	if !contains(args, "-r") {
		t.Error("Recursive=true should produce -r")
	}
}

// TestToArgs_SkipExistingFalse verifies that SkipExisting=false produces -f.
func TestToArgs_SkipExistingFalse(t *testing.T) {
	opts := DefaultConvertOptions()
	opts.SkipExisting = false
	args := opts.ToArgs()

	if !contains(args, "-f") {
		t.Error("SkipExisting=false should produce -f")
	}
}

// TestToArgs_SkipExistingTrue verifies that SkipExisting=true omits -f.
func TestToArgs_SkipExistingTrue(t *testing.T) {
	opts := DefaultConvertOptions()
	opts.SkipExisting = true
	args := opts.ToArgs()

	if contains(args, "-f") {
		t.Error("SkipExisting=true should NOT produce -f")
	}
}

// TestToArgs_Uncompressed verifies that Compression="uncompressed" produces -c uncompressed.
func TestToArgs_Uncompressed(t *testing.T) {
	opts := DefaultConvertOptions()
	opts.Compression = "uncompressed"
	args := opts.ToArgs()

	if !containsSeq(args, "-c", "uncompressed") {
		t.Error("Compression='uncompressed' should produce '-c uncompressed'")
	}
}

// TestToArgs_CropActivearea verifies that Crop="activearea" produces --crop activearea.
func TestToArgs_CropActivearea(t *testing.T) {
	opts := DefaultConvertOptions()
	opts.Crop = "activearea"
	args := opts.ToArgs()

	if !containsSeq(args, "--crop", "activearea") {
		t.Error("Crop='activearea' should produce '--crop activearea'")
	}
}

// TestToArgs_ArtistSet verifies that a non-empty Artist produces --artist <value>.
func TestToArgs_ArtistSet(t *testing.T) {
	opts := DefaultConvertOptions()
	opts.Artist = "Jane Smith"
	args := opts.ToArgs()

	if !containsSeq(args, "--artist", "Jane Smith") {
		t.Error("Artist='Jane Smith' should produce '--artist Jane Smith'")
	}
}

// TestToArgs_ArtistEmpty verifies that an empty Artist omits --artist.
func TestToArgs_ArtistEmpty(t *testing.T) {
	opts := DefaultConvertOptions()
	opts.Artist = ""
	args := opts.ToArgs()

	if contains(args, "--artist") {
		t.Error("Artist='' should NOT produce --artist")
	}
}
