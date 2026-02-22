package converter

import "testing"

func TestPDFToHEIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "pdf", "heif")

	c := NewPDFToHEIFConverter()
	input := mustEncodeFormat(t, "pdf")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "heif")
}

func TestPDFToHEIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewPDFToHEIFConverter()
	assertInvalidInputError(t, c, "pdf", "heif")
}
