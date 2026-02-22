package converter

import "testing"

func TestJPEGToWEBPConverterConvert(t *testing.T) {
	requireFormatPairSupport(t, "jpeg", "webp")

	c := NewJPEGToWEBPConverter()
	input := mustEncodeFormat(t, "jpeg")

	output, err := c.Convert(input)
	if err != nil {
		t.Fatalf("expected conversion to succeed, got error: %v", err)
	}

	assertOutputFormat(t, output, "webp")
}

func TestJPEGToWEBPConverterConvertRejectsInvalidInput(t *testing.T) {
	c := NewJPEGToWEBPConverter()
	assertInvalidInputError(t, c, "jpeg", "webp")
}
