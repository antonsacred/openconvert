package converter

import "testing"

func TestPDFToWEBPConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "pdf", "webp")

	c := NewPDFToWEBPConverter()
	input := mustEncodeFormat(t, "pdf")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "webp")
}

func TestPDFToWEBPConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewPDFToWEBPConverter()
	assertInvalidInputError(t, c, "pdf", "webp")
}
