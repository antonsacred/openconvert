package converter

import "testing"

func TestWEBPToJPGConverterConvert(t *testing.T) {
	c := NewWEBPToJPGConverter()
	input := mustEncodeWEBP(t)

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertJPEGOutput(t, output)
}

func TestWEBPToJPGConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewWEBPToJPGConverter()
	assertInvalidInputError(t, c, "webp", "jpg")
}
