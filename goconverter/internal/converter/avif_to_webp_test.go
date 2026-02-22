package converter

import "testing"

func TestAVIFToWEBPConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "avif", "webp")

	c := NewAVIFToWEBPConverter()
	input := mustEncodeFormat(t, "avif")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "webp")
}

func TestAVIFToWEBPConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewAVIFToWEBPConverter()
	assertInvalidInputError(t, c, "avif", "webp")
}
