package converter

import "testing"

func TestMAGICKToWEBPConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "magick", "webp")

	c := NewMAGICKToWEBPConverter()
	input := mustEncodeFormat(t, "magick")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "webp")
}

func TestMAGICKToWEBPConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewMAGICKToWEBPConverter()
	assertInvalidInputError(t, c, "magick", "webp")
}
