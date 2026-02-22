package converter

import "testing"

func TestAVIFToJPEGConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "avif", "jpeg")

	c := NewAVIFToJPEGConverter()
	input := mustEncodeFormat(t, "avif")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "jpeg")
}

func TestAVIFToJPEGConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewAVIFToJPEGConverter()
	assertInvalidInputError(t, c, "avif", "jpeg")
}
