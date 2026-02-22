package converter

import "testing"

func TestHEIFToWEBPConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "heif", "webp")

	c := NewHEIFToWEBPConverter()
	input := mustEncodeFormat(t, "heif")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "webp")
}

func TestHEIFToWEBPConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewHEIFToWEBPConverter()
	assertInvalidInputError(t, c, "heif", "webp")
}
