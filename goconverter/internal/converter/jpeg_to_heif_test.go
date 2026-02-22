package converter

import "testing"

func TestJPEGToHEIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "jpeg", "heif")

	c := NewJPEGToHEIFConverter()
	input := mustEncodeFormat(t, "jpeg")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "heif")
}

func TestJPEGToHEIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewJPEGToHEIFConverter()
	assertInvalidInputError(t, c, "jpeg", "heif")
}
