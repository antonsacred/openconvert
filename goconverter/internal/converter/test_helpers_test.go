package converter

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"strings"
	"testing"

	"github.com/h2non/bimg"
)

func mustEncodePNG(t *testing.T) []byte {
	t.Helper()

	img := fixtureImage()
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("failed to encode png fixture: %v", err)
	}

	return buf.Bytes()
}

func mustEncodeJPG(t *testing.T) []byte {
	t.Helper()

	img := fixtureImage()
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatalf("failed to encode jpg fixture: %v", err)
	}

	return buf.Bytes()
}

func mustEncodeWEBP(t *testing.T) []byte {
	t.Helper()

	pngInput := mustEncodePNG(t)
	output, err := bimg.NewImage(pngInput).Convert(bimg.WEBP)
	if err != nil {
		t.Fatalf("failed to encode webp fixture: %v", err)
	}

	return output
}

func assertJPEGOutput(t *testing.T, output []byte) {
	t.Helper()

	decoded, err := jpeg.Decode(bytes.NewReader(output))
	if err != nil {
		t.Fatalf("expected valid jpeg output, got error: %v", err)
	}
	assertImageSize(t, decoded)
}

func assertPNGOutput(t *testing.T, output []byte) {
	t.Helper()

	decoded, err := png.Decode(bytes.NewReader(output))
	if err != nil {
		t.Fatalf("expected valid png output, got error: %v", err)
	}
	assertImageSize(t, decoded)
}

func assertWEBPOutput(t *testing.T, output []byte) {
	t.Helper()

	if got := bimg.DetermineImageTypeName(output); got != "webp" {
		t.Fatalf("expected webp output, got %q", got)
	}

	size, err := bimg.Size(output)
	if err != nil {
		t.Fatalf("failed to read webp output size: %v", err)
	}
	if size.Width != 2 || size.Height != 2 {
		t.Fatalf("expected output size 2x2, got %dx%d", size.Width, size.Height)
	}
}

func assertInvalidInputError(t *testing.T, c Converter, source string, target string) {
	t.Helper()

	_, err := c.Convert([]byte("invalid-input-data"))
	if err == nil {
		t.Fatalf("expected conversion to fail for invalid input data")
	}

	expected := "convert " + source + " to " + target
	if !strings.Contains(err.Error(), expected) {
		t.Fatalf("expected wrapped conversion error %q, got: %v", expected, err)
	}
}

func fixtureImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.RGBA{R: 255, G: 0, B: 0, A: 255}), image.Point{}, draw.Src)
	img.Set(1, 0, color.RGBA{R: 0, G: 255, B: 0, A: 255})
	img.Set(0, 1, color.RGBA{R: 0, G: 0, B: 255, A: 255})
	img.Set(1, 1, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	return img
}

func assertImageSize(t *testing.T, img image.Image) {
	t.Helper()

	if img.Bounds().Dx() != 2 || img.Bounds().Dy() != 2 {
		t.Fatalf("expected output size 2x2, got %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
	}
}
