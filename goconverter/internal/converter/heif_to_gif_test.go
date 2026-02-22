package converter

import "testing"

func TestHEIFToGIFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "heif", "gif")

	c := NewHEIFToGIFConverter()
	input := mustEncodeFormat(t, "heif")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "gif")
}

func TestHEIFToGIFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewHEIFToGIFConverter()
	assertInvalidInputError(t, c, "heif", "gif")
}
