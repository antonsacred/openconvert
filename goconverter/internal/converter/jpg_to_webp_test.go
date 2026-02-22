package converter

import "testing"

func TestJPGToWEBPConverterConvert(t *testing.T) {
	c := NewJPGToWEBPConverter()
	input := mustEncodeJPG(t)

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertWEBPOutput(t, output)
}

func TestJPGToWEBPConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewJPGToWEBPConverter()
	assertInvalidInputError(t, c, "jpg", "webp")
}
