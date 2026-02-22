package converter

import "testing"

func TestPNGToWEBPConverterConvert(t *testing.T) {
	c := NewPNGToWEBPConverter()
	input := mustEncodePNG(t)

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertWEBPOutput(t, output)
}

func TestPNGToWEBPConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewPNGToWEBPConverter()
	assertInvalidInputError(t, c, "png", "webp")
}
