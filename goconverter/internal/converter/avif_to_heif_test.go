package converter

import "testing"

func TestAVIFToHEIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "avif", "heif")

	c := NewAVIFToHEIFConverter()
	input := mustEncodeFormat(t, "avif")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "heif")
}

func TestAVIFToHEIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewAVIFToHEIFConverter()
	assertInvalidInputError(t, c, "avif", "heif")
}
