package converter

import "github.com/h2non/bimg"

type PNGToJPEGConverter struct{}

var _ Converter = (*PNGToJPEGConverter)(nil)

func NewPNGToJPEGConverter() *PNGToJPEGConverter {
	return &PNGToJPEGConverter{}
}

func (c *PNGToJPEGConverter) SourceFormat() string {
	return "png"
}

func (c *PNGToJPEGConverter) TargetFormat() string {
	return "jpeg"
}

func (c *PNGToJPEGConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.JPEG)
}
