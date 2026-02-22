package converter

import "testing"

func TestMAGICKToGIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "magick", "gif")

	c := NewMAGICKToGIFConverter()
	input := mustEncodeFormat(t, "magick")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "gif")
}

func TestMAGICKToGIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewMAGICKToGIFConverter()
	assertInvalidInputError(t, c, "magick", "gif")
}
