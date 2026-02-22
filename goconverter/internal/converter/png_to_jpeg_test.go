package converter

import "testing"

func TestPNGToJPEGConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "png", "jpeg")

	c := NewPNGToJPEGConverter()
	input := mustEncodeFormat(t, "png")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "jpeg")
}

func TestPNGToJPEGConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewPNGToJPEGConverter()
	assertInvalidInputError(t, c, "png", "jpeg")
}
