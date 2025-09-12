//go:build !windows

package agnostic

import (
	"bytes"
	"image"

	"github.com/chai2010/webp"
)

// encodeWebP uses native WebP encoder (requires libwebp)
func encodeWebP(img image.Image) (*bytes.Buffer, error) {
	var webpBuf bytes.Buffer
	err := webp.Encode(&webpBuf, img, &webp.Options{Lossless: true})
	if err != nil {
		return nil, err
	}
	return &webpBuf, nil
}
