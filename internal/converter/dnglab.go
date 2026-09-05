package converter

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DNGLab wraps the dnglab CLI binary.
type DNGLab struct {
	BinaryPath string
}

// DetectBinary locates the dnglab binary. It first checks the directory that
// contains the running executable, then falls back to the system PATH.
func (d *DNGLab) DetectBinary() error {
	// (1) Same directory as the running executable
	execPath, err := os.Executable()
	if err == nil {
		candidate := filepath.Join(filepath.Dir(execPath), "dnglab")
		if _, statErr := os.Stat(candidate); statErr == nil {
			d.BinaryPath = candidate
			return nil
		}
	}

	// (2) System PATH
	path, err := exec.LookPath("dnglab")
	if err == nil {
		d.BinaryPath = path
		return nil
	}

	return fmt.Errorf("dnglab binary not found: not in the same directory as the application and not on PATH")
}

// Version runs 'dnglab --version' and returns the trimmed output.
func (d *DNGLab) Version() (string, error) {
	if d.BinaryPath == "" {
		return "", fmt.Errorf("dnglab binary path is not set; call DetectBinary first")
	}

	var out bytes.Buffer
	cmd := exec.Command(d.BinaryPath, "--version") //nolint:gosec
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("running dnglab --version: %w", err)
	}

	return strings.TrimSpace(out.String()), nil
}

// Convert converts all RAW files found under opts.InputPath to DNG format.
// progress is called before each file with the file path, 1-based index, and total count.
// All errors are collected and returned as a single combined error.
func (d *DNGLab) Convert(ctx context.Context, opts ConvertOptions, progress func(file string, current, total int)) error {
	if d.BinaryPath == "" {
		return fmt.Errorf("dnglab binary path is not set; call DetectBinary first")
	}

	files, err := ScanForRawFiles(opts.InputPath, opts.Recursive)
	if err != nil {
		return fmt.Errorf("scanning for RAW files: %w", err)
	}

	total := len(files)
	var errs []string

	for i, file := range files {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if progress != nil {
			progress(file, i+1, total)
		}

		// Build argument list: flags first, then positional args
		args := []string{"convert"}
		args = append(args, opts.ToArgs()...)
		args = append(args, file)
		if opts.OutputPath != "" {
			args = append(args, opts.OutputPath)
		}

		var out bytes.Buffer
		cmd := exec.CommandContext(ctx, d.BinaryPath, args...) //nolint:gosec
		cmd.Stdout = &out
		cmd.Stderr = &out

		if runErr := cmd.Run(); runErr != nil {
			errs = append(errs, fmt.Sprintf("converting %s: %v (output: %s)", file, runErr, strings.TrimSpace(out.String())))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%d file(s) failed to convert:\n%s", len(errs), strings.Join(errs, "\n"))
	}

	return nil
}
