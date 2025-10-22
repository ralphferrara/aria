package webp

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"os/exec"
	"runtime"
	"strconv"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp" // WebP decoder
)

// ResizeToWebPMax resizes an image to fit within (maxW x maxH) and encodes WebP via cwebp.
// - quality: 1..100 (typical: 80-85)
// - returns WebP bytes
func ResizeToWebPMax(src []byte, maxW, maxH, quality int) ([]byte, error) {
	if maxW <= 0 || maxH <= 0 {
		return nil, fmt.Errorf("maxW/maxH must be > 0")
	}
	if quality < 1 || quality > 100 {
		quality = 82
	}

	// 1) Decode any common raster format
	img, _, err := image.Decode(bytes.NewReader(src))
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	// 2) Compute target size (fit within, never upscale)
	b := img.Bounds()
	sw, sh := b.Dx(), b.Dy()
	nw, nh := fitWithin(sw, sh, maxW, maxH)
	if nw == sw && nh == sh {
		// still pass through for encoding, but skip resize
	} else {
		// 3) High-quality resize (CatmullRom)
		dst := image.NewRGBA(image.Rect(0, 0, nw, nh))
		draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
		img = dst
	}

	// 4) Encode to an intermediate PNG (lossless) for cwebp stdin
	var pngBuf bytes.Buffer
	if err := png.Encode(&pngBuf, img); err != nil {
		return nil, fmt.Errorf("encode png: %w", err)
	}

	// 5) Call cwebp (fast, robust encoder)
	cwebp := "cwebp"
	if runtime.GOOS == "windows" {
		cwebp = "cwebp.exe"
	}
	cmd := exec.Command(cwebp, "-q", strconv.Itoa(quality), "-quiet", "--mt", "-", "-o", "-")
	cmd.Stdin = &pngBuf
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("cwebp: %v (%s)", err, errBuf.String())
	}
	return out.Bytes(), nil
}

// fitWithin returns new size that fits inside (mw x mh), keeping aspect, never upscaling.
func fitWithin(w, h, mw, mh int) (int, int) {
	sf := math.Min(float64(mw)/float64(w), float64(mh)/float64(h))
	if sf > 1 {
		sf = 1 // no upscaling
	}
	nw := int(math.Round(float64(w) * sf))
	nh := int(math.Round(float64(h) * sf))
	if nw < 1 {
		nw = 1
	}
	if nh < 1 {
		nh = 1
	}
	return nw, nh
}

// Optional: helpers for strict JPEG/PNG decode if you want to avoid the generic image.Decode.
var _ = jpeg.DefaultQuality
var _ = gif.DisposalBackground
