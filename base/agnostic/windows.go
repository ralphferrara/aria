//go:build windows

package agnostic

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var cwebpPath = `C:\webp\bin\cwebp.exe` // adjust if installed elsewhere

func encodeWebP(img image.Image, lossless bool, quality int) (*bytes.Buffer, error) {
	if quality < 1 || quality > 100 {
		quality = 82
	}

	// Encode PNG to memory
	var pngBuf bytes.Buffer
	if err := png.Encode(&pngBuf, img); err != nil {
		return nil, err
	}

	// First attempt: stdin, using `-- -` to mark "-" as the input file
	buf, err := runCwebpStdin(&pngBuf, lossless, quality)
	if err == nil {
		return buf, nil
	}
	// If stdin is not supported, cwebp often prints "Unknown option '-'"
	// or similar; fall back to temp-file path invocation.
	if !isStdinUnsupported(err) {
		// some other error; return it
		return nil, err
	}

	// Fallback: write temp PNG, call cwebp with file path
	return runCwebpTempFile(&pngBuf, lossless, quality)
}

func runCwebpStdin(pngData io.Reader, lossless bool, quality int) (*bytes.Buffer, error) {
	args := []string{"-quiet", "--"}
	// quality/lossless first
	if lossless {
		args = append([]string{"-lossless"}, args...)
	} else {
		args = append([]string{"-q", strconv.Itoa(quality)}, args...)
	}
	// input from stdin, output to stdout
	args = append(args, "-", "-o", "-")

	cmd := exec.Command(cwebpPath, args...)
	cmd.Stdin = readerToReadSeeker(pngData) // ensure it's seekable for some shells
	var out bytes.Buffer
	cmd.Stdout = &out
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("cwebp(stdin) failed: %w | %s", err, strings.TrimSpace(errBuf.String()))
	}
	return &out, nil
}

func runCwebpTempFile(pngData io.Reader, lossless bool, quality int) (*bytes.Buffer, error) {
	dir := os.TempDir()
	inFile, err := os.CreateTemp(dir, "img-*.png")
	if err != nil {
		return nil, err
	}
	defer os.Remove(inFile.Name())
	defer inFile.Close()

	if _, err := io.Copy(inFile, pngData); err != nil {
		return nil, err
	}
	if err := inFile.Close(); err != nil {
		return nil, err
	}

	// Use stdout for output to avoid creating another temp file
	args := []string{"-quiet"}
	if lossless {
		args = append(args, "-lossless")
	} else {
		args = append(args, "-q", strconv.Itoa(quality))
	}
	args = append(args, inFile.Name(), "-o", "-")

	cmd := exec.Command(cwebpPath, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("cwebp(file) failed: %w | %s", err, strings.TrimSpace(errBuf.String()))
	}

	return &out, nil
}

// Some Windows builds treat "-" as an unknown option and print a message.
// Detect that to decide when to fall back to temp files.
func isStdinUnsupported(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "Unknown option '-'") ||
		strings.Contains(msg, "Unknown option '--'") ||
		strings.Contains(strings.ToLower(msg), "unknown option") ||
		strings.Contains(strings.ToLower(msg), "invalid argument")
}

// In case some shells need a ReadSeeker on Stdin, wrap the reader into a buffer.
func readerToReadSeeker(r io.Reader) io.ReadSeeker {
	switch v := r.(type) {
	case io.ReadSeeker:
		return v
	default:
		var b bytes.Buffer
		_, _ = io.Copy(&b, r)
		return bytes.NewReader(b.Bytes())
	}
}
