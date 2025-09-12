package agnostic

import (
	"bytes"
	"image"
)

// EncodeWebP encodes an image.Image to WebP format and returns a buffer or error.
func EncodeWebP(img image.Image) (*bytes.Buffer, error) {
	return encodeWebP(img)
}
