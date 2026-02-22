package converter

import "testing"

func TestWEBPToTIFFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "webp", "tiff")

	c := NewWEBPToTIFFConverter()
	input := mustEncodeFormat(t, "webp")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "tiff")
}

func TestWEBPToTIFFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewWEBPToTIFFConverter()
	assertInvalidInputError(t, c, "webp", "tiff")
}
