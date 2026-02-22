package converter

import "github.com/h2non/bimg"

type JPEGToWEBPConverter struct{}

var _ Converter = (*JPEGToWEBPConverter)(nil)

func NewJPEGToWEBPConverter() *JPEGToWEBPConverter {
	return &JPEGToWEBPConverter{}
}

func (c *JPEGToWEBPConverter) SourceFormat() string {
	return "jpeg"
}

func (c *JPEGToWEBPConverter) TargetFormat() string {
	return "webp"
}

func (c *JPEGToWEBPConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.WEBP)
}
