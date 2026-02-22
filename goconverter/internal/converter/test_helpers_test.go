package converter

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/h2non/bimg"
)

var (
	bimgModuleDirOnce sync.Once
	bimgModuleDir     string
	bimgModuleDirErr  error

	bimgTestdataCacheMu sync.Mutex
	bimgTestdataCache   = map[string][]byte{}
)

func requireFormatPairSupport(t *testing.T, source string, target string) {
	t.Helper()

	if !bimg.IsTypeNameSupported(source) {
		t.Skipf("source format %q is not load-supported by current libvips build", source)
	}
	if !bimg.IsTypeNameSupportedSave(target) {
		t.Skipf("target format %q is not save-supported by current libvips build", target)
	}
}

func mustEncodeFormat(t *testing.T, format string) []byte {
	t.Helper()

	switch format {
	case "avif":
		return mustEncodeWithBIMG(t, "avif", bimg.AVIF)
	case "gif":
		return mustEncodeWithBIMG(t, "gif", bimg.GIF)
	case "heif":
		return mustEncodeWithBIMG(t, "heif", bimg.HEIF)
	case "jpeg":
		return mustEncodeJPEG(t)
	case "magick":
		// Converter selection is request-driven; conversion does not validate
		// source bytes against SourceFormat(). Use a stable decodable fixture.
		return mustEncodePNG(t)
	case "pdf":
		return mustReadBIMGTestdataFile(t, "test.pdf")
	case "png":
		return mustEncodePNG(t)
	case "svg":
		return mustEncodeSVG()
	case "tiff":
		return mustEncodeWithBIMG(t, "tiff", bimg.TIFF)
	case "webp":
		return mustEncodeWithBIMG(t, "webp", bimg.WEBP)
	default:
		t.Fatalf("unsupported source fixture format %q", format)
		return nil
	}
}

func mustEncodeWithBIMG(t *testing.T, format string, imageType bimg.ImageType) []byte {
	t.Helper()

	if !bimg.IsTypeNameSupportedSave(format) {
		t.Skipf("cannot generate %q source fixture: format is not save-supported", format)
	}

	output, err := bimg.NewImage(mustEncodePNG(t)).Convert(imageType)
	if err != nil {
		t.Fatalf("failed to encode %s fixture: %v", format, err)
	}

	if got := bimg.DetermineImageTypeName(output); got != format {
		t.Fatalf("expected %s fixture, got %q", format, got)
	}

	return output
}

func mustEncodePNG(t *testing.T) []byte {
	t.Helper()

	img := fixtureImage()
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("failed to encode png fixture: %v", err)
	}

	return buf.Bytes()
}

func mustEncodeJPEG(t *testing.T) []byte {
	t.Helper()

	img := fixtureImage()
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatalf("failed to encode jpeg fixture: %v", err)
	}

	return buf.Bytes()
}

func mustEncodeSVG() []byte {
	return []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="2" height="2"><rect width="2" height="2" fill="red"/></svg>`)
}

func assertOutputFormat(t *testing.T, output []byte, expectedFormat string) {
	t.Helper()

	if got := bimg.DetermineImageTypeName(output); got != expectedFormat {
		t.Fatalf("expected %s output, got %q", expectedFormat, got)
	}

	size, err := bimg.Size(output)
	if err != nil {
		t.Fatalf("failed to read output size for %s: %v", expectedFormat, err)
	}
	if size.Width <= 0 || size.Height <= 0 {
		t.Fatalf("expected positive output dimensions, got %dx%d", size.Width, size.Height)
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

func mustReadBIMGTestdataFile(t *testing.T, fileName string) []byte {
	t.Helper()

	bimgTestdataCacheMu.Lock()
	if fixture, ok := bimgTestdataCache[fileName]; ok {
		output := make([]byte, len(fixture))
		copy(output, fixture)
		bimgTestdataCacheMu.Unlock()
		return output
	}
	bimgTestdataCacheMu.Unlock()

	moduleDir, err := bimgModuleDirPath()
	if err != nil {
		t.Skipf("unable to locate bimg module directory: %v", err)
	}

	fixturePath := filepath.Join(moduleDir, "testdata", fileName)
	fixture, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Skipf("unable to read bimg test fixture %q: %v", fixturePath, err)
	}

	bimgTestdataCacheMu.Lock()
	bimgTestdataCache[fileName] = fixture
	bimgTestdataCacheMu.Unlock()

	output := make([]byte, len(fixture))
	copy(output, fixture)
	return output
}

func bimgModuleDirPath() (string, error) {
	bimgModuleDirOnce.Do(func() {
		cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/h2non/bimg")
		output, err := cmd.CombinedOutput()
		if err != nil {
			bimgModuleDirErr = fmt.Errorf("resolve bimg module directory: %w (output: %s)", err, strings.TrimSpace(string(output)))
			return
		}

		moduleDir := strings.TrimSpace(string(output))
		if moduleDir == "" {
			bimgModuleDirErr = fmt.Errorf("resolve bimg module directory: empty output")
			return
		}

		bimgModuleDir = moduleDir
	})

	if bimgModuleDirErr != nil {
		return "", bimgModuleDirErr
	}

	return bimgModuleDir, nil
}
