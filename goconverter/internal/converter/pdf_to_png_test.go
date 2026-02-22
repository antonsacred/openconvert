package converter

import "testing"

func TestPDFToPNGConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "pdf", "png")

	c := NewPDFToPNGConverter()
	input := mustEncodeFormat(t, "pdf")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "png")
}

func TestPDFToPNGConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewPDFToPNGConverter()
	assertInvalidInputError(t, c, "pdf", "png")
}
