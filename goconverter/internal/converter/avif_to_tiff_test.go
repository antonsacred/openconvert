package converter

import "testing"

func TestAVIFToTIFFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "avif", "tiff")

	c := NewAVIFToTIFFConverter()
	input := mustEncodeFormat(t, "avif")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "tiff")
}

func TestAVIFToTIFFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewAVIFToTIFFConverter()
	assertInvalidInputError(t, c, "avif", "tiff")
}
