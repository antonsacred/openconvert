package converter

import "testing"

func TestWEBPToHEIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "webp", "heif")

	c := NewWEBPToHEIFConverter()
	input := mustEncodeFormat(t, "webp")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "heif")
}

func TestWEBPToHEIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewWEBPToHEIFConverter()
	assertInvalidInputError(t, c, "webp", "heif")
}
