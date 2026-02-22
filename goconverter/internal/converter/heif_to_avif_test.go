package converter

import "testing"

func TestHEIFToAVIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "heif", "avif")

	c := NewHEIFToAVIFConverter()
	input := mustEncodeFormat(t, "heif")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "avif")
}

func TestHEIFToAVIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewHEIFToAVIFConverter()
	assertInvalidInputError(t, c, "heif", "avif")
}
