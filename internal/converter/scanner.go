package converter

import (
	"os"
	"path/filepath"
	"strings"
)

// supportedExtensions is the set of RAW file extensions that dnglab can convert.
// Keys are stored in uppercase for case-insensitive matching.
var supportedExtensions = map[string]struct{}{
	".CR3": {},
	".CR2": {},
	".CRW": {},
	".NEF": {},
	".NRW": {},
	".ARW": {},
	".SRF": {},
	".SR2": {},
	".RAF": {},
	".RW2": {},
	".ORF": {},
	".PEF": {},
	".3FR": {},
	".ARI": {},
	".ERF": {},
	".KDC": {},
	".DCS": {},
	".DCR": {},
	".IIQ": {},
	".MOS": {},
	".MEF": {},
	".MRW": {},
	".SRW": {},
}

// IsSupportedRawFile reports whether filename has a supported RAW extension.
// The check is case-insensitive.
func IsSupportedRawFile(filename string) bool {
	ext := strings.ToUpper(filepath.Ext(filename))
	if ext == "" {
		return false
	}
	_, ok := supportedExtensions[ext]
	return ok
}

// ScanForRawFiles walks path and returns the absolute paths of all RAW files found.
// When recursive is false only the top-level directory is searched.
func ScanForRawFiles(path string, recursive bool) ([]string, error) {
	var files []string

	err := filepath.WalkDir(path, func(entry string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip subdirectories when not in recursive mode
		if d.IsDir() {
			if entry != path && !recursive {
				return filepath.SkipDir
			}
			return nil
		}

		if IsSupportedRawFile(d.Name()) {
			abs, absErr := filepath.Abs(entry)
			if absErr != nil {
				return absErr
			}
			files = append(files, abs)
		}

		return nil
	})

	return files, err
}
