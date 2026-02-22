package converter

import "github.com/h2non/bimg"

type TIFFToGIFConverter struct{}

var _ Converter = (*TIFFToGIFConverter)(nil)

func NewTIFFToGIFConverter() *TIFFToGIFConverter {
	return &TIFFToGIFConverter{}
}

func (c *TIFFToGIFConverter) SourceFormat() string {
	return "tiff"
}

func (c *TIFFToGIFConverter) TargetFormat() string {
	return "gif"
}

func (c *TIFFToGIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.GIF)
}
