package converter

import "github.com/h2non/bimg"

type GIFToPNGConverter struct{}

var _ Converter = (*GIFToPNGConverter)(nil)

func NewGIFToPNGConverter() *GIFToPNGConverter {
	return &GIFToPNGConverter{}
}

func (c *GIFToPNGConverter) SourceFormat() string {
	return "gif"
}

func (c *GIFToPNGConverter) TargetFormat() string {
	return "png"
}

func (c *GIFToPNGConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.PNG)
}
