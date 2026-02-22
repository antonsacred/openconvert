package converter

import "github.com/h2non/bimg"

type TIFFToHEIFConverter struct{}

var _ Converter = (*TIFFToHEIFConverter)(nil)

func NewTIFFToHEIFConverter() *TIFFToHEIFConverter {
	return &TIFFToHEIFConverter{}
}

func (c *TIFFToHEIFConverter) SourceFormat() string {
	return "tiff"
}

func (c *TIFFToHEIFConverter) TargetFormat() string {
	return "heif"
}

func (c *TIFFToHEIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.HEIF)
}
