package converter

import "testing"

func TestSVGToPNGConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "svg", "png")

	c := NewSVGToPNGConverter()
	input := mustEncodeFormat(t, "svg")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "png")
}

func TestSVGToPNGConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewSVGToPNGConverter()
	assertInvalidInputError(t, c, "svg", "png")
}
