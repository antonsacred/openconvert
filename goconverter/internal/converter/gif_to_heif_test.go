package converter

import "testing"

func TestGIFToHEIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "gif", "heif")

	c := NewGIFToHEIFConverter()
	input := mustEncodeFormat(t, "gif")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "heif")
}

func TestGIFToHEIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewGIFToHEIFConverter()
	assertInvalidInputError(t, c, "gif", "heif")
}
