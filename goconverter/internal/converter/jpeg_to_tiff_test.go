package converter

import "testing"

func TestJPEGToTIFFConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "jpeg", "tiff")

	c := NewJPEGToTIFFConverter()
	input := mustEncodeFormat(t, "jpeg")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "tiff")
}

func TestJPEGToTIFFConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewJPEGToTIFFConverter()
	assertInvalidInputError(t, c, "jpeg", "tiff")
}
