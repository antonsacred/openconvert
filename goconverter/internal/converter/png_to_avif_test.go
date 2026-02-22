package converter

import "testing"

func TestPNGToAVIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "png", "avif")

	c := NewPNGToAVIFConverter()
	input := mustEncodeFormat(t, "png")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "avif")
}

func TestPNGToAVIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewPNGToAVIFConverter()
	assertInvalidInputError(t, c, "png", "avif")
}
