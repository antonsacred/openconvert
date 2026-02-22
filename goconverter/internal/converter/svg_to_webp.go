package converter

import "github.com/h2non/bimg"

type SVGToWEBPConverter struct{}

var _ Converter = (*SVGToWEBPConverter)(nil)

func NewSVGToWEBPConverter() *SVGToWEBPConverter {
	return &SVGToWEBPConverter{}
}

func (c *SVGToWEBPConverter) SourceFormat() string {
	return "svg"
}

func (c *SVGToWEBPConverter) TargetFormat() string {
	return "webp"
}

func (c *SVGToWEBPConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.WEBP)
}
