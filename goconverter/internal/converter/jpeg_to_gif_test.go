package converter

import "testing"

func TestJPEGToGIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "jpeg", "gif")

	c := NewJPEGToGIFConverter()
	input := mustEncodeFormat(t, "jpeg")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "gif")
}

func TestJPEGToGIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewJPEGToGIFConverter()
	assertInvalidInputError(t, c, "jpeg", "gif")
}
