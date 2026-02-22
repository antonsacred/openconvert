package converter

import "github.com/h2non/bimg"

type SVGToAVIFConverter struct{}

var _ Converter = (*SVGToAVIFConverter)(nil)

func NewSVGToAVIFConverter() *SVGToAVIFConverter {
	return &SVGToAVIFConverter{}
}

func (c *SVGToAVIFConverter) SourceFormat() string {
	return "svg"
}

func (c *SVGToAVIFConverter) TargetFormat() string {
	return "avif"
}

func (c *SVGToAVIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.AVIF)
}
