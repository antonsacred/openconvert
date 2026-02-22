package converter

import "github.com/h2non/bimg"

type JPEGToAVIFConverter struct{}

var _ Converter = (*JPEGToAVIFConverter)(nil)

func NewJPEGToAVIFConverter() *JPEGToAVIFConverter {
	return &JPEGToAVIFConverter{}
}

func (c *JPEGToAVIFConverter) SourceFormat() string {
	return "jpeg"
}

func (c *JPEGToAVIFConverter) TargetFormat() string {
	return "avif"
}

func (c *JPEGToAVIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.AVIF)
}
