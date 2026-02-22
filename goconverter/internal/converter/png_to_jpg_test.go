package converter

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"strings"
	"testing"
)

func TestPNGToJPGConverterConvert(t *testing.T) {
	c := NewPNGToJPGConverter()
	input := mustEncodePNG(t)

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	decoded, err := jpeg.Decode(bytes.NewReader(output))
	if err != nil {
		t.Fatalf("expected valid jpeg output, got error: %v", err)
	}

	if decoded.Bounds().Dx() != 2 || decoded.Bounds().Dy() != 2 {
		t.Fatalf("expected output size 2x2, got %dx%d", decoded.Bounds().Dx(), decoded.Bounds().Dy())
	}
}

func TestPNGToJPGConverterConvertRejectsInvalidPNG(t *testing.T) {
	c := NewPNGToJPGConverter()

	_, err := c.Convert([]byte("invalid-png-data"))
	if err == nil {
		t.Fatalf("expected conversion to fail for invalid png data")
	}
	if !strings.Contains(err.Error(), "decode png") {
		t.Fatalf("expected decode png error, got: %v", err)
	}
}

func TestPossibleConversions(t *testing.T) {
	output := PossibleConversions()

	if len(output) != 1 {
		t.Fatalf("expected exactly 1 registered conversion, got %d", len(output))
	}

	if output[0][0] != "png" || output[0][1] != "jpg" {
		t.Fatalf("expected conversion pair [png, jpg], got [%s, %s]", output[0][0], output[0][1])
	}
}

func mustEncodePNG(t *testing.T) []byte {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	img.Set(1, 0, color.RGBA{R: 0, G: 255, B: 0, A: 255})
	img.Set(0, 1, color.RGBA{R: 0, G: 0, B: 255, A: 255})
	img.Set(1, 1, color.RGBA{R: 255, G: 255, B: 255, A: 255})

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("failed to encode png fixture: %v", err)
	}

	return buf.Bytes()
}
