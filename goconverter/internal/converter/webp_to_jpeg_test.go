package converter

import "testing"

func TestWEBPToJPEGConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "webp", "jpeg")

	c := NewWEBPToJPEGConverter()
	input := mustEncodeFormat(t, "webp")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "jpeg")
}

func TestWEBPToJPEGConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewWEBPToJPEGConverter()
	assertInvalidInputError(t, c, "webp", "jpeg")
}
