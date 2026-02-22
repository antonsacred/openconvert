package converter

import "testing"

func TestPNGToWEBPConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "png", "webp")

	c := NewPNGToWEBPConverter()
	input := mustEncodeFormat(t, "png")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "webp")
}

func TestPNGToWEBPConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewPNGToWEBPConverter()
	assertInvalidInputError(t, c, "png", "webp")
}
