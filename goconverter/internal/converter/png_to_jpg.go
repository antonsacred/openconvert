package converter

import "github.com/h2non/bimg"

type PNGToJPGConverter struct{}

var _ Converter = (*PNGToJPGConverter)(nil)

func NewPNGToJPGConverter() *PNGToJPGConverter {
	return &PNGToJPGConverter{}
}

func (c *PNGToJPGConverter) SourceFormat() string {
	return "png"
}

func (c *PNGToJPGConverter) TargetFormat() string {
	return "jpg"
}

func (c *PNGToJPGConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.JPEG)
}
