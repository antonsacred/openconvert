package converter

import "testing"

func TestPDFToAVIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "pdf", "avif")

	c := NewPDFToAVIFConverter()
	input := mustEncodeFormat(t, "pdf")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "avif")
}

func TestPDFToAVIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewPDFToAVIFConverter()
	assertInvalidInputError(t, c, "pdf", "avif")
}
