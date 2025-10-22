package agnostic

import (
	"bytes"
	"image"
)

// EncodeWebP encodes an image.Image to WebP.
// If lossless == true, quality is ignored by some encoders (use 100).
// quality: 1..100 (typical lossy: 80-85)
func EncodeWebP(img image.Image, lossless bool, quality int) (*bytes.Buffer, error) {
	return encodeWebP(img, lossless, quality)
}
