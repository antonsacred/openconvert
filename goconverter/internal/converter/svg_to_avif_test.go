package converter

import "testing"

func TestSVGToAVIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "svg", "avif")

	c := NewSVGToAVIFConverter()
	input := mustEncodeFormat(t, "svg")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "avif")
}

func TestSVGToAVIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewSVGToAVIFConverter()
	assertInvalidInputError(t, c, "svg", "avif")
}
