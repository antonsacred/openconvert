package converter

import "github.com/h2non/bimg"

type GIFToJPEGConverter struct{}

var _ Converter = (*GIFToJPEGConverter)(nil)

func NewGIFToJPEGConverter() *GIFToJPEGConverter {
	return &GIFToJPEGConverter{}
}

func (c *GIFToJPEGConverter) SourceFormat() string {
	return "gif"
}

func (c *GIFToJPEGConverter) TargetFormat() string {
	return "jpeg"
}

func (c *GIFToJPEGConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.JPEG)
}
