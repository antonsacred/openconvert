package converter

import "testing"

func TestSVGToTIFFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "svg", "tiff")

	c := NewSVGToTIFFConverter()
	input := mustEncodeFormat(t, "svg")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "tiff")
}

func TestSVGToTIFFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewSVGToTIFFConverter()
	assertInvalidInputError(t, c, "svg", "tiff")
}
