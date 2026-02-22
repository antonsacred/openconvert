package converter

import "github.com/h2non/bimg"

type GIFToAVIFConverter struct{}

var _ Converter = (*GIFToAVIFConverter)(nil)

func NewGIFToAVIFConverter() *GIFToAVIFConverter {
	return &GIFToAVIFConverter{}
}

func (c *GIFToAVIFConverter) SourceFormat() string {
	return "gif"
}

func (c *GIFToAVIFConverter) TargetFormat() string {
	return "avif"
}

func (c *GIFToAVIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.AVIF)
}
