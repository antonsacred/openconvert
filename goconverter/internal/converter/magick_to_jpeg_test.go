package converter

import "testing"

func TestMAGICKToJPEGConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "magick", "jpeg")

	c := NewMAGICKToJPEGConverter()
	input := mustEncodeFormat(t, "magick")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "jpeg")
}

func TestMAGICKToJPEGConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewMAGICKToJPEGConverter()
	assertInvalidInputError(t, c, "magick", "jpeg")
}
