package converter

import "fmt"

// ConvertOptions holds all settings for a dnglab conversion run.
type ConvertOptions struct {
	InputPath    string
	OutputPath   string // empty = same location as input
	Recursive    bool
	SkipExisting bool   // GUI "skip if exists" -- true means do NOT pass -f/--override
	Compression  string // "lossless" or "uncompressed"
	Crop         string // "best", "activearea", or "none"
	EmbedRaw     bool
	Preview      bool
	Thumbnail    bool
	Artist       string
}

// DefaultConvertOptions returns sensible defaults for a conversion run.
func DefaultConvertOptions() ConvertOptions {
	return ConvertOptions{
		Compression:  "lossless",
		Crop:         "best",
		EmbedRaw:     true,
		Preview:      true,
		Thumbnail:    true,
		SkipExisting: true,
	}
}

// ToArgs builds the dnglab CLI flags for this set of options.
// It does NOT include the positional INPUT and OUTPUT arguments --
// those are appended by the caller.
func (o ConvertOptions) ToArgs() []string {
	var args []string

	if o.Recursive {
		args = append(args, "-r")
	}

	// -f means "override existing files"; SkipExisting=true means omit -f
	if !o.SkipExisting {
		args = append(args, "-f")
	}

	args = append(args, "-c", o.Compression)
	args = append(args, "--crop", o.Crop)
	args = append(args, "--embed-raw", fmt.Sprintf("%v", o.EmbedRaw))
	args = append(args, "--dng-preview", fmt.Sprintf("%v", o.Preview))
	args = append(args, "--dng-thumbnail", fmt.Sprintf("%v", o.Thumbnail))

	if o.Artist != "" {
		args = append(args, "--artist", o.Artist)
	}

	return args
}
