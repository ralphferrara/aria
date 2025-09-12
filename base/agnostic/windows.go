//go:build windows

package agnostic

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
)

// encodeWebP uses cwebp on Windows
func encodeWebP(img image.Image) (*bytes.Buffer, error) {
	var pngBuf bytes.Buffer
	if err := png.Encode(&pngBuf, img); err != nil {
		return nil, err
	}

	cmd := exec.Command("C:\\webp\\bin\\cwebp", "-lossless", "-quiet", "-o", "-", "--", "-")
	cmd.Stdin = &pngBuf

	var webpBuf bytes.Buffer
	cmd.Stdout = &webpBuf
	cmd.Stderr = os.Stderr // <-- ADD THIS to show real error

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("WebP encoding failed: %w", err)
	}

	return &webpBuf, nil
}
