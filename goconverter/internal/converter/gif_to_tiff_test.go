package converter

import "testing"

func TestGIFToTIFFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "gif", "tiff")

	c := NewGIFToTIFFConverter()
	input := mustEncodeFormat(t, "gif")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "tiff")
}

func TestGIFToTIFFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewGIFToTIFFConverter()
	assertInvalidInputError(t, c, "gif", "tiff")
}
