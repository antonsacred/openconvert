package converter

import "github.com/h2non/bimg"

type SVGToPNGConverter struct{}

var _ Converter = (*SVGToPNGConverter)(nil)

func NewSVGToPNGConverter() *SVGToPNGConverter {
	return &SVGToPNGConverter{}
}

func (c *SVGToPNGConverter) SourceFormat() string {
	return "svg"
}

func (c *SVGToPNGConverter) TargetFormat() string {
	return "png"
}

func (c *SVGToPNGConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.PNG)
}
