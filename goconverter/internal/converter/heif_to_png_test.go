package converter

import "testing"

func TestHEIFToPNGConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "heif", "png")

	c := NewHEIFToPNGConverter()
	input := mustEncodeFormat(t, "heif")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "png")
}

func TestHEIFToPNGConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewHEIFToPNGConverter()
	assertInvalidInputError(t, c, "heif", "png")
}
