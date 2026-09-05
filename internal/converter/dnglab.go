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

// Convert converts all RAW files under opts.InputPath to DNG format by
// invoking dnglab with the input directory directly. This lets dnglab use its
// own parallelism and avoids per-file process-spawn overhead.
//
// totalFiles is the number of RAW files the caller already scanned (passed
// through to the progress callback so the UI can show a meaningful count).
// The caller is responsible for scanning upfront; Convert does not scan again.
// Context cancellation is propagated directly into the subprocess.
func (d *DNGLab) Convert(ctx context.Context, opts ConvertOptions, totalFiles int, progress func(file string, current, total int)) error {
	if d.BinaryPath == "" {
		return fmt.Errorf("dnglab binary path is not set; call DetectBinary first")
	}

	if progress != nil {
		progress("", 0, totalFiles)
	}

	// Build argument list: flags first, then the input directory (and optional
	// output directory). dnglab processes the whole directory in one shot,
	// which preserves its internal parallelism.
	args := []string{"convert"}
	args = append(args, opts.ToArgs()...)
	args = append(args, opts.InputPath)
	if opts.OutputPath != "" {
		args = append(args, opts.OutputPath)
	}

	var out bytes.Buffer
	cmd := exec.CommandContext(ctx, d.BinaryPath, args...) //nolint:gosec
	cmd.Stdout = &out
	cmd.Stderr = &out

	if runErr := cmd.Run(); runErr != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("dnglab convert failed: %v\noutput:\n%s", runErr, strings.TrimSpace(out.String()))
	}

	return nil
}
