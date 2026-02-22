package converter

import "testing"

func TestHEIFToTIFFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "heif", "tiff")

	c := NewHEIFToTIFFConverter()
	input := mustEncodeFormat(t, "heif")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "tiff")
}

func TestHEIFToTIFFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewHEIFToTIFFConverter()
	assertInvalidInputError(t, c, "heif", "tiff")
}
