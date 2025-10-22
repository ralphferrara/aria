//go:build windows

package agnostic

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"strconv"
)

var cwebpPath = `C:\webp\bin\cwebp.exe` // change if installed elsewhere

func encodeWebP(img image.Image, lossless bool, quality int) (*bytes.Buffer, error) {
	if quality < 1 || quality > 100 {
		quality = 82
	}
	// Encode PNG to feed cwebp stdin
	var pngBuf bytes.Buffer
	if err := png.Encode(&pngBuf, img); err != nil {
		return nil, err
	}

	args := []string{"-quiet", "--mt", "-", "-o", "-"}
	if lossless {
		args = append([]string{"-lossless"}, args...)
	} else {
		args = append([]string{"-q", strconv.Itoa(quality)}, args...)
	}

	cmd := exec.Command(cwebpPath, args...)
	cmd.Stdin = &pngBuf

	var webpBuf bytes.Buffer
	cmd.Stdout = &webpBuf
	cmd.Stderr = os.Stderr // keep this to see real encoder errors

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("WebP encoding failed: %w", err)
	}
	return &webpBuf, nil
}
