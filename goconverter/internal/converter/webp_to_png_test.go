package converter

import "testing"

func TestWEBPToPNGConverterConvert(t *testing.T) {
	c := NewWEBPToPNGConverter()
	input := mustEncodeWEBP(t)

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertPNGOutput(t, output)
}

func TestWEBPToPNGConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewWEBPToPNGConverter()
	assertInvalidInputError(t, c, "webp", "png")
}
