package converter

import "testing"

func TestJPEGToAVIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "jpeg", "avif")

	c := NewJPEGToAVIFConverter()
	input := mustEncodeFormat(t, "jpeg")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "avif")
}

func TestJPEGToAVIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewJPEGToAVIFConverter()
	assertInvalidInputError(t, c, "jpeg", "avif")
}
