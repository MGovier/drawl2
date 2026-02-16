package ai

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"strings"

	"golang.org/x/image/draw"
)

const maxVisionDim = 512

// resizeDataURL takes a data:image/png;base64,... string, decodes it,
// resizes if larger than maxVisionDim, and returns a new data URL.
func resizeDataURL(dataURL string) (string, error) {
	b64 := strings.TrimPrefix(dataURL, "data:image/png;base64,")
	if b64 == dataURL {
		// Not a png data URL, pass through
		return dataURL, nil
	}

	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", fmt.Errorf("decode base64: %w", err)
	}

	src, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return "", fmt.Errorf("decode image: %w", err)
	}

	bounds := src.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	// Only resize if larger than limit
	if w <= maxVisionDim && h <= maxVisionDim {
		return dataURL, nil
	}

	// Scale preserving aspect ratio
	newW, newH := maxVisionDim, maxVisionDim
	if w > h {
		newH = h * maxVisionDim / w
	} else {
		newW = w * maxVisionDim / h
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.BiLinear.Scale(dst, dst.Bounds(), src, bounds, draw.Over, nil)

	var buf bytes.Buffer
	if err := png.Encode(&buf, dst); err != nil {
		return "", fmt.Errorf("encode png: %w", err)
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
