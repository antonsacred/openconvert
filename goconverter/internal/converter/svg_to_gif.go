package converter

import "github.com/h2non/bimg"

type SVGToGIFConverter struct{}

var _ Converter = (*SVGToGIFConverter)(nil)

func NewSVGToGIFConverter() *SVGToGIFConverter {
	return &SVGToGIFConverter{}
}

func (c *SVGToGIFConverter) SourceFormat() string {
	return "svg"
}

func (c *SVGToGIFConverter) TargetFormat() string {
	return "gif"
}

func (c *SVGToGIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.GIF)
}
