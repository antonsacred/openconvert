package converter

import "testing"

func TestTIFFToGIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "tiff", "gif")

	c := NewTIFFToGIFConverter()
	input := mustEncodeFormat(t, "tiff")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "gif")
}

func TestTIFFToGIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewTIFFToGIFConverter()
	assertInvalidInputError(t, c, "tiff", "gif")
}
