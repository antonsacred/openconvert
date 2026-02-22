package converter

import "testing"

func TestMAGICKToAVIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "magick", "avif")

	c := NewMAGICKToAVIFConverter()
	input := mustEncodeFormat(t, "magick")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "avif")
}

func TestMAGICKToAVIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewMAGICKToAVIFConverter()
	assertInvalidInputError(t, c, "magick", "avif")
}
