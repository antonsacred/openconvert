package converter

import "testing"

func TestMAGICKToTIFFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "magick", "tiff")

	c := NewMAGICKToTIFFConverter()
	input := mustEncodeFormat(t, "magick")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "tiff")
}

func TestMAGICKToTIFFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewMAGICKToTIFFConverter()
	assertInvalidInputError(t, c, "magick", "tiff")
}
