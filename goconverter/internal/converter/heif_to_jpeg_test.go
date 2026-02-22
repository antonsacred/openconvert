package converter

import "testing"

func TestHEIFToJPEGConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "heif", "jpeg")

	c := NewHEIFToJPEGConverter()
	input := mustEncodeFormat(t, "heif")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "jpeg")
}

func TestHEIFToJPEGConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewHEIFToJPEGConverter()
	assertInvalidInputError(t, c, "heif", "jpeg")
}
