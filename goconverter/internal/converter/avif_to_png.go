package converter

import "github.com/h2non/bimg"

type AVIFToPNGConverter struct{}

var _ Converter = (*AVIFToPNGConverter)(nil)

func NewAVIFToPNGConverter() *AVIFToPNGConverter {
	return &AVIFToPNGConverter{}
}

func (c *AVIFToPNGConverter) SourceFormat() string {
	return "avif"
}

func (c *AVIFToPNGConverter) TargetFormat() string {
	return "png"
}

func (c *AVIFToPNGConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.PNG)
}
