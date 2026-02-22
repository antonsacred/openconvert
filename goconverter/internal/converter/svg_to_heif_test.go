package converter

import "testing"

func TestSVGToHEIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "svg", "heif")

	c := NewSVGToHEIFConverter()
	input := mustEncodeFormat(t, "svg")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "heif")
}

func TestSVGToHEIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewSVGToHEIFConverter()
	assertInvalidInputError(t, c, "svg", "heif")
}
