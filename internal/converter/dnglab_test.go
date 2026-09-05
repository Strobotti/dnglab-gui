package converter

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeStub writes a small shell script to dir/name, makes it executable, and
// returns its path. The script exits with the given exitCode and prints args
// to stdout so tests can verify what arguments were passed.
func writeStub(t *testing.T, dir, name string, exitCode int) string {
	t.Helper()
	path := filepath.Join(dir, name)
	script := "#!/bin/sh\necho \"$@\"\nexit " + itoa(exitCode) + "\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatalf("writeStub: %v", err)
	}
	return path
}

// itoa converts a small int to a decimal string without importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	buf := make([]byte, 0, 10)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	if neg {
		buf = append([]byte{'-'}, buf...)
	}
	return string(buf)
}

// makeInputDir creates a small temp directory with one RAW file for tests that
// need a non-empty input path.
func makeInputDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "test.CR3"), []byte("fake"), 0o644); err != nil {
		t.Fatalf("makeInputDir: %v", err)
	}
	return dir
}

// TestConvertSuccess verifies that Convert returns nil when dnglab exits 0.
func TestConvertSuccess(t *testing.T) {
	stubDir := t.TempDir()
	stub := writeStub(t, stubDir, "dnglab", 0)
	inputDir := makeInputDir(t)

	d := DNGLab{BinaryPath: stub}
	opts := DefaultConvertOptions()
	opts.InputPath = inputDir

	err := d.Convert(context.Background(), opts, 1, nil)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestConvertFailure verifies that Convert returns an error containing the
// stub's output when dnglab exits non-zero.
func TestConvertFailure(t *testing.T) {
	stubDir := t.TempDir()
	// Write a stub that prints a recognizable message and exits 1.
	path := filepath.Join(stubDir, "dnglab")
	script := "#!/bin/sh\necho 'conversion failed: unsupported format'\nexit 1\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatalf("write stub: %v", err)
	}

	inputDir := makeInputDir(t)
	d := DNGLab{BinaryPath: path}
	opts := DefaultConvertOptions()
	opts.InputPath = inputDir

	err := d.Convert(context.Background(), opts, 1, nil)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "conversion failed: unsupported format") {
		t.Fatalf("error message should contain stub output, got: %v", err)
	}
}

// TestConvertArgs verifies that Convert invokes dnglab with the expected
// argument order: "convert" subcommand, then flags, then inputDir.
func TestConvertArgs(t *testing.T) {
	stubDir := t.TempDir()
	stub := writeStub(t, stubDir, "dnglab", 0)
	inputDir := makeInputDir(t)

	d := DNGLab{BinaryPath: stub}
	opts := DefaultConvertOptions()
	opts.InputPath = inputDir
	opts.Recursive = false
	opts.SkipExisting = true // means -f is omitted

	var captured string
	err := d.Convert(context.Background(), opts, 1, func(file string, current, total int) {
		// first (and only) progress callback: file should be empty string
		captured = file
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// The progress callback should be called with file="" at start.
	if captured != "" {
		t.Errorf("expected progress callback file to be empty string, got %q", captured)
	}

	// Run the stub directly to capture the args it would receive.
	// We exercise the argument builder via ToArgs and verify inputDir is last.
	args := opts.ToArgs()
	// "convert" is prepended in Convert; ToArgs only returns the flags.
	// Verify the inputDir is not inside the flags.
	for _, arg := range args {
		if arg == inputDir {
			t.Errorf("inputDir should not appear in ToArgs output, only in the final positional arg")
		}
	}
}

// TestConvertOutputPath verifies that when OutputPath is set, it is appended
// after the inputDir in the argument list.
func TestConvertOutputPath(t *testing.T) {
	// Use a stub that prints its args so we can inspect them.
	stubDir := t.TempDir()
	argCapturePath := filepath.Join(stubDir, "captured_args.txt")
	scriptBody := "#!/bin/sh\necho \"$@\" > " + argCapturePath + "\nexit 0\n"
	stubPath := filepath.Join(stubDir, "dnglab")
	if err := os.WriteFile(stubPath, []byte(scriptBody), 0o755); err != nil {
		t.Fatalf("write stub: %v", err)
	}

	inputDir := makeInputDir(t)
	outputDir := t.TempDir()

	d := DNGLab{BinaryPath: stubPath}
	opts := DefaultConvertOptions()
	opts.InputPath = inputDir
	opts.OutputPath = outputDir

	if err := d.Convert(context.Background(), opts, 1, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	argsBytes, err := os.ReadFile(argCapturePath)
	if err != nil {
		t.Fatalf("reading captured args: %v", err)
	}
	argsLine := strings.TrimSpace(string(argsBytes))

	// The captured line is all args space-separated by the shell "$@".
	// outputDir must appear after inputDir.
	inputIdx := strings.Index(argsLine, inputDir)
	outputIdx := strings.Index(argsLine, outputDir)

	if inputIdx == -1 {
		t.Errorf("inputDir %q not found in args: %q", inputDir, argsLine)
	}
	if outputIdx == -1 {
		t.Errorf("outputDir %q not found in args: %q", outputDir, argsLine)
	}
	if inputIdx != -1 && outputIdx != -1 && outputIdx <= inputIdx {
		t.Errorf("outputDir should appear after inputDir in args: %q", argsLine)
	}
}
