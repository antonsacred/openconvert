package converter

import "github.com/h2non/bimg"

type PNGToTIFFConverter struct{}

var _ Converter = (*PNGToTIFFConverter)(nil)

func NewPNGToTIFFConverter() *PNGToTIFFConverter {
	return &PNGToTIFFConverter{}
}

func (c *PNGToTIFFConverter) SourceFormat() string {
	return "png"
}

func (c *PNGToTIFFConverter) TargetFormat() string {
	return "tiff"
}

func (c *PNGToTIFFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.TIFF)
}
