package converter

import "testing"

func TestJPEGToPNGConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "jpeg", "png")

	c := NewJPEGToPNGConverter()
	input := mustEncodeFormat(t, "jpeg")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "png")
}

func TestJPEGToPNGConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewJPEGToPNGConverter()
	assertInvalidInputError(t, c, "jpeg", "png")
}
