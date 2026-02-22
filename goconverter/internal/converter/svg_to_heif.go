package converter

import "github.com/h2non/bimg"

type SVGToHEIFConverter struct{}

var _ Converter = (*SVGToHEIFConverter)(nil)

func NewSVGToHEIFConverter() *SVGToHEIFConverter {
	return &SVGToHEIFConverter{}
}

func (c *SVGToHEIFConverter) SourceFormat() string {
	return "svg"
}

func (c *SVGToHEIFConverter) TargetFormat() string {
	return "heif"
}

func (c *SVGToHEIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.HEIF)
}
