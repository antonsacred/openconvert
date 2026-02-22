package converter

import "github.com/h2non/bimg"

type SVGToJPEGConverter struct{}

var _ Converter = (*SVGToJPEGConverter)(nil)

func NewSVGToJPEGConverter() *SVGToJPEGConverter {
	return &SVGToJPEGConverter{}
}

func (c *SVGToJPEGConverter) SourceFormat() string {
	return "svg"
}

func (c *SVGToJPEGConverter) TargetFormat() string {
	return "jpeg"
}

func (c *SVGToJPEGConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.JPEG)
}
