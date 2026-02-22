package converter

import "github.com/h2non/bimg"

type PNGToGIFConverter struct{}

var _ Converter = (*PNGToGIFConverter)(nil)

func NewPNGToGIFConverter() *PNGToGIFConverter {
	return &PNGToGIFConverter{}
}

func (c *PNGToGIFConverter) SourceFormat() string {
	return "png"
}

func (c *PNGToGIFConverter) TargetFormat() string {
	return "gif"
}

func (c *PNGToGIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.GIF)
}
