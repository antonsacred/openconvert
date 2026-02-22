package converter

import "testing"

func TestSVGToGIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "svg", "gif")

	c := NewSVGToGIFConverter()
	input := mustEncodeFormat(t, "svg")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "gif")
}

func TestSVGToGIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewSVGToGIFConverter()
	assertInvalidInputError(t, c, "svg", "gif")
}
