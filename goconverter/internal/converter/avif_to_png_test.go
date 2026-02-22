package converter

import "testing"

func TestAVIFToPNGConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "avif", "png")

	c := NewAVIFToPNGConverter()
	input := mustEncodeFormat(t, "avif")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "png")
}

func TestAVIFToPNGConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewAVIFToPNGConverter()
	assertInvalidInputError(t, c, "avif", "png")
}
