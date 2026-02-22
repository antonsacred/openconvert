package converter

import "testing"

func TestMAGICKToHEIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "magick", "heif")

	c := NewMAGICKToHEIFConverter()
	input := mustEncodeFormat(t, "magick")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "heif")
}

func TestMAGICKToHEIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewMAGICKToHEIFConverter()
	assertInvalidInputError(t, c, "magick", "heif")
}
