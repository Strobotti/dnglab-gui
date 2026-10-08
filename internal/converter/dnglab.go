package converter

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
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

// Result summarizes a conversion run.
type Result struct {
	Converted int
	Skipped   int          // destination already existed and was not overwritten
	Failed    []FailedFile // files dnglab could not convert
}

// FailedFile describes one file that could not be converted.
type FailedFile struct {
	Path   string
	Reason string
}

// Convert converts all RAW files under opts.InputPath to DNG format by
// invoking dnglab with the input directory directly, which keeps dnglab's own
// parallelism. The output of "dnglab -v" is parsed live and reported through
// progress. totalFiles is the number of RAW files the caller already scanned.
//
// A non-zero exit status from dnglab is not an error when it converted,
// skipped or failed individual files; dnglab reports partial failures that way.
// Result.Failed carries those files. An error is returned only when dnglab
// cannot run, the run was cancelled, or it failed without processing any file.
func (d *DNGLab) Convert(ctx context.Context, opts ConvertOptions, totalFiles int, progress func(Event)) (Result, error) {
	var res Result
	if d.BinaryPath == "" {
		return res, fmt.Errorf("dnglab binary path is not set; call DetectBinary first")
	}

	args := []string{"convert", "-v"}
	args = append(args, opts.ToArgs()...)
	args = append(args, opts.InputPath)
	if opts.OutputPath != "" {
		args = append(args, opts.OutputPath)
	}

	pr, pw := io.Pipe()
	cmd := exec.CommandContext(ctx, d.BinaryPath, args...) //nolint:gosec
	cmd.Stdout = pw
	cmd.Stderr = pw

	var (
		mu  sync.Mutex
		out strings.Builder
	)
	done := make(chan struct{})
	go func() {
		defer close(done)
		sc := bufio.NewScanner(pr)
		sc.Buffer(make([]byte, 64*1024), 1024*1024)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" {
				continue
			}
			mu.Lock()
			out.WriteString(line + "\n")
			mu.Unlock()

			ev := parseLine(line)
			switch {
			case ev.Converted:
				res.Converted++
			case ev.Skipped:
				res.Skipped++
			case ev.Failed:
				res.Failed = append(res.Failed, FailedFile{Path: ev.Source, Reason: ev.Line})
			}
			if progress != nil {
				progress(ev)
			}
		}
		_, _ = io.Copy(io.Discard, pr)
	}()

	runErr := cmd.Run()
	_ = pw.Close()
	<-done

	if ctx.Err() != nil {
		return res, ctx.Err()
	}
	if runErr != nil {
		processed := res.Converted + res.Skipped + len(res.Failed)
		var exitErr *exec.ExitError
		if errors.As(runErr, &exitErr) && processed > 0 {
			return res, nil
		}
		mu.Lock()
		output := strings.TrimSpace(out.String())
		mu.Unlock()
		return res, fmt.Errorf("dnglab convert failed: %v\noutput:\n%s", runErr, output)
	}

	return res, nil
}

// Event is one parsed line of dnglab's verbose output.
type Event struct {
	Line      string // raw output line
	Source    string // source file path, for converted, skipped and failed files
	Converted bool   // dnglab reported a successfully converted file
	Skipped   bool   // destination already exists and the file was not overwritten
	Failed    bool   // dnglab could not convert the file
}

// parseLine classifies a line of dnglab's -v output. Known per-file formats:
//
//	Status: Converted '/path/in.ARW' => '/path/out.dng' (in 0.93s)
//	Status: Failed: '/path/in.ARW', Already exists: /path/out.dng
//
// Other lines (summaries, "Error: ..." repeats) are only logged, not counted.
func parseLine(line string) Event {
	ev := Event{Line: line}
	switch {
	case strings.HasPrefix(line, "Status: Converted"):
		ev.Converted = true
		ev.Source = quotedValue(line)
	case strings.HasPrefix(line, "Status: Failed") && strings.Contains(line, "Already exists"):
		ev.Skipped = true
		ev.Source = quotedValue(line)
	case strings.HasPrefix(line, "Status: Failed"):
		ev.Failed = true
		ev.Source = quotedValue(line)
	}
	return ev
}

// quotedValue returns the text between the first pair of single quotes in s.
func quotedValue(s string) string {
	i := strings.Index(s, "'")
	if i < 0 {
		return ""
	}
	rest := s[i+1:]
	if j := strings.Index(rest, "'"); j >= 0 {
		return rest[:j]
	}
	return ""
}
