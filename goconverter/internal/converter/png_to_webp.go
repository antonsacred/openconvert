package converter

import "github.com/h2non/bimg"

type PNGToWEBPConverter struct{}

var _ Converter = (*PNGToWEBPConverter)(nil)

func NewPNGToWEBPConverter() *PNGToWEBPConverter {
	return &PNGToWEBPConverter{}
}

func (c *PNGToWEBPConverter) SourceFormat() string {
	return "png"
}

func (c *PNGToWEBPConverter) TargetFormat() string {
	return "webp"
}

func (c *PNGToWEBPConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.WEBP)
}
