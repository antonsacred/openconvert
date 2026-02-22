package converter

import "testing"

func TestPDFToTIFFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "pdf", "tiff")

	c := NewPDFToTIFFConverter()
	input := mustEncodeFormat(t, "pdf")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "tiff")
}

func TestPDFToTIFFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewPDFToTIFFConverter()
	assertInvalidInputError(t, c, "pdf", "tiff")
}
