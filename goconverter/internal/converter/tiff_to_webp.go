package converter

import "github.com/h2non/bimg"

type TIFFToWEBPConverter struct{}

var _ Converter = (*TIFFToWEBPConverter)(nil)

func NewTIFFToWEBPConverter() *TIFFToWEBPConverter {
	return &TIFFToWEBPConverter{}
}

func (c *TIFFToWEBPConverter) SourceFormat() string {
	return "tiff"
}

func (c *TIFFToWEBPConverter) TargetFormat() string {
	return "webp"
}

func (c *TIFFToWEBPConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.WEBP)
}
