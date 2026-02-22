package converter

import "testing"

func TestGIFToAVIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "gif", "avif")

	c := NewGIFToAVIFConverter()
	input := mustEncodeFormat(t, "gif")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "avif")
}

func TestGIFToAVIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewGIFToAVIFConverter()
	assertInvalidInputError(t, c, "gif", "avif")
}
