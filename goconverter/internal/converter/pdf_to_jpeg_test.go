package converter

import "testing"

func TestPDFToJPEGConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "pdf", "jpeg")

	c := NewPDFToJPEGConverter()
	input := mustEncodeFormat(t, "pdf")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "jpeg")
}

func TestPDFToJPEGConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewPDFToJPEGConverter()
	assertInvalidInputError(t, c, "pdf", "jpeg")
}
