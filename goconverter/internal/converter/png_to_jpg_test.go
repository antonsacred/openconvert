package converter

import "testing"

func TestPNGToJPGConverterConvert(t *testing.T) {
	c := NewPNGToJPGConverter()
	input := mustEncodePNG(t)

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertJPEGOutput(t, output)
}

func TestPNGToJPGConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewPNGToJPGConverter()
	assertInvalidInputError(t, c, "png", "jpg")
}
