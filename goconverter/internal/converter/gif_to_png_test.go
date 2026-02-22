package converter

import "testing"

func TestGIFToPNGConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "gif", "png")

	c := NewGIFToPNGConverter()
	input := mustEncodeFormat(t, "gif")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "png")
}

func TestGIFToPNGConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewGIFToPNGConverter()
	assertInvalidInputError(t, c, "gif", "png")
}
