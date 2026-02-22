package converter

import "github.com/h2non/bimg"

type SVGToTIFFConverter struct{}

var _ Converter = (*SVGToTIFFConverter)(nil)

func NewSVGToTIFFConverter() *SVGToTIFFConverter {
	return &SVGToTIFFConverter{}
}

func (c *SVGToTIFFConverter) SourceFormat() string {
	return "svg"
}

func (c *SVGToTIFFConverter) TargetFormat() string {
	return "tiff"
}

func (c *SVGToTIFFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.TIFF)
}
