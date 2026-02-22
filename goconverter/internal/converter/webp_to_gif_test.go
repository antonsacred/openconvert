package converter

import "testing"

func TestWEBPToGIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "webp", "gif")

	c := NewWEBPToGIFConverter()
	input := mustEncodeFormat(t, "webp")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "gif")
}

func TestWEBPToGIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewWEBPToGIFConverter()
	assertInvalidInputError(t, c, "webp", "gif")
}
