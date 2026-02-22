package converter

import "testing"

func TestPNGToHEIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "png", "heif")

	c := NewPNGToHEIFConverter()
	input := mustEncodeFormat(t, "png")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "heif")
}

func TestPNGToHEIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewPNGToHEIFConverter()
	assertInvalidInputError(t, c, "png", "heif")
}
