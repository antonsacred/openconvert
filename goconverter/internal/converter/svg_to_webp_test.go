package converter

import "testing"

func TestSVGToWEBPConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "svg", "webp")

	c := NewSVGToWEBPConverter()
	input := mustEncodeFormat(t, "svg")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "webp")
}

func TestSVGToWEBPConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewSVGToWEBPConverter()
	assertInvalidInputError(t, c, "svg", "webp")
}
