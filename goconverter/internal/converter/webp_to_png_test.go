package converter

import "testing"

func TestWEBPToPNGConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "webp", "png")

	c := NewWEBPToPNGConverter()
	input := mustEncodeFormat(t, "webp")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "png")
}

func TestWEBPToPNGConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewWEBPToPNGConverter()
	assertInvalidInputError(t, c, "webp", "png")
}
