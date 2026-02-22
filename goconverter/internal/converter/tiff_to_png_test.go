package converter

import "testing"

func TestTIFFToPNGConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "tiff", "png")

	c := NewTIFFToPNGConverter()
	input := mustEncodeFormat(t, "tiff")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "png")
}

func TestTIFFToPNGConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewTIFFToPNGConverter()
	assertInvalidInputError(t, c, "tiff", "png")
}
