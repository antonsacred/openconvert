package converter

import "testing"

func TestPDFToGIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "pdf", "gif")

	c := NewPDFToGIFConverter()
	input := mustEncodeFormat(t, "pdf")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "gif")
}

func TestPDFToGIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewPDFToGIFConverter()
	assertInvalidInputError(t, c, "pdf", "gif")
}
