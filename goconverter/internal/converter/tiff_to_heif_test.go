package converter

import "testing"

func TestTIFFToHEIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "tiff", "heif")

	c := NewTIFFToHEIFConverter()
	input := mustEncodeFormat(t, "tiff")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "heif")
}

func TestTIFFToHEIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewTIFFToHEIFConverter()
	assertInvalidInputError(t, c, "tiff", "heif")
}
