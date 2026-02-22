package converter

import "github.com/h2non/bimg"

type GIFToWEBPConverter struct{}

var _ Converter = (*GIFToWEBPConverter)(nil)

func NewGIFToWEBPConverter() *GIFToWEBPConverter {
	return &GIFToWEBPConverter{}
}

func (c *GIFToWEBPConverter) SourceFormat() string {
	return "gif"
}

func (c *GIFToWEBPConverter) TargetFormat() string {
	return "webp"
}

func (c *GIFToWEBPConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.WEBP)
}
