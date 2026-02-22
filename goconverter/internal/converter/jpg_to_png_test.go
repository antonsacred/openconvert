package converter

import "testing"

func TestJPGToPNGConverterConvert(t *testing.T) {
	c := NewJPGToPNGConverter()
	input := mustEncodeJPG(t)

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertPNGOutput(t, output)
}

func TestJPGToPNGConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewJPGToPNGConverter()
	assertInvalidInputError(t, c, "jpg", "png")
}
