package converter

import "testing"

func TestTIFFToJPEGConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "tiff", "jpeg")

	c := NewTIFFToJPEGConverter()
	input := mustEncodeFormat(t, "tiff")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "jpeg")
}

func TestTIFFToJPEGConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewTIFFToJPEGConverter()
	assertInvalidInputError(t, c, "tiff", "jpeg")
}
