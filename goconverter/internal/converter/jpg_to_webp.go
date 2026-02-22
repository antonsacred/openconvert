package converter

import "github.com/h2non/bimg"

type JPGToWEBPConverter struct{}

var _ Converter = (*JPGToWEBPConverter)(nil)

func NewJPGToWEBPConverter() *JPGToWEBPConverter {
	return &JPGToWEBPConverter{}
}

func (c *JPGToWEBPConverter) SourceFormat() string {
	return "jpg"
}

func (c *JPGToWEBPConverter) TargetFormat() string {
	return "webp"
}

func (c *JPGToWEBPConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.WEBP)
}
