package converter

import "testing"

func TestGIFToWEBPConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "gif", "webp")

	c := NewGIFToWEBPConverter()
	input := mustEncodeFormat(t, "gif")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "webp")
}

func TestGIFToWEBPConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewGIFToWEBPConverter()
	assertInvalidInputError(t, c, "gif", "webp")
}
