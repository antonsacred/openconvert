package converter

import "testing"

func TestTIFFToWEBPConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "tiff", "webp")

	c := NewTIFFToWEBPConverter()
	input := mustEncodeFormat(t, "tiff")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "webp")
}

func TestTIFFToWEBPConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewTIFFToWEBPConverter()
	assertInvalidInputError(t, c, "tiff", "webp")
}
