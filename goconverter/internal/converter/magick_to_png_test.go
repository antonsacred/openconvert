package converter

import "testing"

func TestMAGICKToPNGConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "magick", "png")

	c := NewMAGICKToPNGConverter()
	input := mustEncodeFormat(t, "magick")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "png")
}

func TestMAGICKToPNGConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewMAGICKToPNGConverter()
	assertInvalidInputError(t, c, "magick", "png")
}
