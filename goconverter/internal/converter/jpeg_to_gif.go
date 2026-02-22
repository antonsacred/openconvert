package converter

import "github.com/h2non/bimg"

type JPEGToGIFConverter struct{}

var _ Converter = (*JPEGToGIFConverter)(nil)

func NewJPEGToGIFConverter() *JPEGToGIFConverter {
	return &JPEGToGIFConverter{}
}

func (c *JPEGToGIFConverter) SourceFormat() string {
	return "jpeg"
}

func (c *JPEGToGIFConverter) TargetFormat() string {
	return "gif"
}

func (c *JPEGToGIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.GIF)
}
