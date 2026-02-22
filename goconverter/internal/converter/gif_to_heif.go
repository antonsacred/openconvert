package converter

import "github.com/h2non/bimg"

type GIFToHEIFConverter struct{}

var _ Converter = (*GIFToHEIFConverter)(nil)

func NewGIFToHEIFConverter() *GIFToHEIFConverter {
	return &GIFToHEIFConverter{}
}

func (c *GIFToHEIFConverter) SourceFormat() string {
	return "gif"
}

func (c *GIFToHEIFConverter) TargetFormat() string {
	return "heif"
}

func (c *GIFToHEIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.HEIF)
}
