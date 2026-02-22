package converter

import "github.com/h2non/bimg"

type PNGToAVIFConverter struct{}

var _ Converter = (*PNGToAVIFConverter)(nil)

func NewPNGToAVIFConverter() *PNGToAVIFConverter {
	return &PNGToAVIFConverter{}
}

func (c *PNGToAVIFConverter) SourceFormat() string {
	return "png"
}

func (c *PNGToAVIFConverter) TargetFormat() string {
	return "avif"
}

func (c *PNGToAVIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.AVIF)
}
