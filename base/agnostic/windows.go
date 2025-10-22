//go:build windows

package agnostic

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

var cwebpPath = `C:\webp\bin\cwebp.exe` // adjust if needed

func encodeWebP(img image.Image, lossless bool, quality int) (*bytes.Buffer, error) {
	if quality < 1 || quality > 100 {
		quality = 82
	}

	// 1) Temp input PNG
	inFile, err := os.CreateTemp("", "in-*.png")
	if err != nil {
		return nil, err
	}
	inPath := inFile.Name()
	defer func() {
		inFile.Close()
		os.Remove(inPath)
	}()

	if err := png.Encode(inFile, img); err != nil {
		return nil, err
	}
	if err := inFile.Close(); err != nil {
		return nil, err
	}

	// 2) Temp output WEBP (explicit file path, no stdout)
	outFile, err := os.CreateTemp("", "out-*.webp")
	if err != nil {
		return nil, err
	}
	outPath := outFile.Name()
	outFile.Close()
	defer os.Remove(outPath)

	// 3) Build args (no stdin/stdout, no --mt)
	args := []string{"-quiet"}
	if lossless {
		args = append(args, "-lossless")
	} else {
		args = append(args, "-q", strconv.Itoa(quality))
	}
	args = append(args, inPath, "-o", outPath)

	cmd := exec.Command(cwebpPath, args...)
	// Optional: see errors
	// cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("WebP encoding failed: %w", err)
	}

	// 4) Read output file
	data, err := os.ReadFile(outPath)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("cwebp produced 0 bytes: %s", filepath.Base(outPath))
	}
	return bytes.NewBuffer(data), nil
}
