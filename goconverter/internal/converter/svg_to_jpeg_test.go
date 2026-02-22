package converter

import "testing"

func TestSVGToJPEGConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "svg", "jpeg")

	c := NewSVGToJPEGConverter()
	input := mustEncodeFormat(t, "svg")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "jpeg")
}

func TestSVGToJPEGConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewSVGToJPEGConverter()
	assertInvalidInputError(t, c, "svg", "jpeg")
}
