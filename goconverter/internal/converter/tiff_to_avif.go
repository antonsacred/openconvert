package converter

import "github.com/h2non/bimg"

type TIFFToAVIFConverter struct{}

var _ Converter = (*TIFFToAVIFConverter)(nil)

func NewTIFFToAVIFConverter() *TIFFToAVIFConverter {
	return &TIFFToAVIFConverter{}
}

func (c *TIFFToAVIFConverter) SourceFormat() string {
	return "tiff"
}

func (c *TIFFToAVIFConverter) TargetFormat() string {
	return "avif"
}

func (c *TIFFToAVIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.AVIF)
}
