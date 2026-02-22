package converter

import "testing"

func TestGIFToJPEGConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "gif", "jpeg")

	c := NewGIFToJPEGConverter()
	input := mustEncodeFormat(t, "gif")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "jpeg")
}

func TestGIFToJPEGConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewGIFToJPEGConverter()
	assertInvalidInputError(t, c, "gif", "jpeg")
}
