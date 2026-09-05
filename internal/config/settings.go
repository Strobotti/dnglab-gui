package config

import "fyne.io/fyne/v2"

const (
	keyCompression  = "compression"
	keyCrop         = "crop"
	keyEmbedRaw     = "embed_raw"
	keyPreview      = "preview"
	keyThumbnail    = "thumbnail"
	keySkipExisting = "skip_existing"
	keyArtist       = "artist"
	keyLastInputDir = "last_input_dir"
	keyLastOutputDir = "last_output_dir"
	keyOutputMode   = "output_mode"
)

// Settings holds all persisted user preferences.
type Settings struct {
	Compression   string
	Crop          string
	EmbedRaw      bool
	Preview       bool
	Thumbnail     bool
	SkipExisting  bool
	Artist        string
	LastInputDir  string
	LastOutputDir string
	OutputMode    string // "same" or "folder"
}

// DefaultSettings returns the default settings, matching DefaultConvertOptions.
func DefaultSettings() Settings {
	return Settings{
		Compression:  "lossless",
		Crop:         "best",
		EmbedRaw:     true,
		Preview:      true,
		Thumbnail:    true,
		SkipExisting: true,
		OutputMode:   "same",
	}
}

// Load reads settings from Fyne preferences, falling back to defaults.
func Load(prefs fyne.Preferences) Settings {
	d := DefaultSettings()
	return Settings{
		Compression:   prefs.StringWithFallback(keyCompression, d.Compression),
		Crop:          prefs.StringWithFallback(keyCrop, d.Crop),
		EmbedRaw:      prefs.BoolWithFallback(keyEmbedRaw, d.EmbedRaw),
		Preview:       prefs.BoolWithFallback(keyPreview, d.Preview),
		Thumbnail:     prefs.BoolWithFallback(keyThumbnail, d.Thumbnail),
		SkipExisting:  prefs.BoolWithFallback(keySkipExisting, d.SkipExisting),
		Artist:        prefs.StringWithFallback(keyArtist, d.Artist),
		LastInputDir:  prefs.StringWithFallback(keyLastInputDir, d.LastInputDir),
		LastOutputDir: prefs.StringWithFallback(keyLastOutputDir, d.LastOutputDir),
		OutputMode:    prefs.StringWithFallback(keyOutputMode, d.OutputMode),
	}
}

// Save writes settings to Fyne preferences.
func Save(prefs fyne.Preferences, s Settings) {
	prefs.SetString(keyCompression, s.Compression)
	prefs.SetString(keyCrop, s.Crop)
	prefs.SetBool(keyEmbedRaw, s.EmbedRaw)
	prefs.SetBool(keyPreview, s.Preview)
	prefs.SetBool(keyThumbnail, s.Thumbnail)
	prefs.SetBool(keySkipExisting, s.SkipExisting)
	prefs.SetString(keyArtist, s.Artist)
	prefs.SetString(keyLastInputDir, s.LastInputDir)
	prefs.SetString(keyLastOutputDir, s.LastOutputDir)
	prefs.SetString(keyOutputMode, s.OutputMode)
}
