//go:build !windows

package agnostic

import (
	"bytes"
	"image"

	"github.com/chai2010/webp"
)

func encodeWebP(img image.Image, lossless bool, quality int) (*bytes.Buffer, error) {
	if quality < 1 || quality > 100 {
		quality = 82
	}
	opts := &webp.Options{
		Lossless: lossless,
		Quality:  float32(quality), // used when Lossless=false
	}
	var webpBuf bytes.Buffer
	if err := webp.Encode(&webpBuf, img, opts); err != nil {
		return nil, err
	}
	return &webpBuf, nil
}
