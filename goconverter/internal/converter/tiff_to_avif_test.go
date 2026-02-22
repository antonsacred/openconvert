package converter

import "testing"

func TestTIFFToAVIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "tiff", "avif")

	c := NewTIFFToAVIFConverter()
	input := mustEncodeFormat(t, "tiff")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "avif")
}

func TestTIFFToAVIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewTIFFToAVIFConverter()
	assertInvalidInputError(t, c, "tiff", "avif")
}
