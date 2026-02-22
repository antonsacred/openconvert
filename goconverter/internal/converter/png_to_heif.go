package converter

import "github.com/h2non/bimg"

type PNGToHEIFConverter struct{}

var _ Converter = (*PNGToHEIFConverter)(nil)

func NewPNGToHEIFConverter() *PNGToHEIFConverter {
	return &PNGToHEIFConverter{}
}

func (c *PNGToHEIFConverter) SourceFormat() string {
	return "png"
}

func (c *PNGToHEIFConverter) TargetFormat() string {
	return "heif"
}

func (c *PNGToHEIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.HEIF)
}
