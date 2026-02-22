package converter

import "testing"

func TestAVIFToGIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "avif", "gif")

	c := NewAVIFToGIFConverter()
	input := mustEncodeFormat(t, "avif")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "gif")
}

func TestAVIFToGIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewAVIFToGIFConverter()
	assertInvalidInputError(t, c, "avif", "gif")
}
