package converter

import "github.com/h2non/bimg"

type JPEGToHEIFConverter struct{}

var _ Converter = (*JPEGToHEIFConverter)(nil)

func NewJPEGToHEIFConverter() *JPEGToHEIFConverter {
	return &JPEGToHEIFConverter{}
}

func (c *JPEGToHEIFConverter) SourceFormat() string {
	return "jpeg"
}

func (c *JPEGToHEIFConverter) TargetFormat() string {
	return "heif"
}

func (c *JPEGToHEIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.HEIF)
}
