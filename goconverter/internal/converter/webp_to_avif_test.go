package converter

import "testing"

func TestWEBPToAVIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "webp", "avif")

	c := NewWEBPToAVIFConverter()
	input := mustEncodeFormat(t, "webp")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "avif")
}

func TestWEBPToAVIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewWEBPToAVIFConverter()
	assertInvalidInputError(t, c, "webp", "avif")
}
