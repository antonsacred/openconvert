package converter

import "testing"

func TestPNGToGIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "png", "gif")

	c := NewPNGToGIFConverter()
	input := mustEncodeFormat(t, "png")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "gif")
}

func TestPNGToGIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewPNGToGIFConverter()
	assertInvalidInputError(t, c, "png", "gif")
}
