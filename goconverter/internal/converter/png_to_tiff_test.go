package converter

import "testing"

func TestPNGToTIFFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "png", "tiff")

	c := NewPNGToTIFFConverter()
	input := mustEncodeFormat(t, "png")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "tiff")
}

func TestPNGToTIFFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewPNGToTIFFConverter()
	assertInvalidInputError(t, c, "png", "tiff")
}
