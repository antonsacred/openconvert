package converter

import "github.com/h2non/bimg"

type MAGICKToHEIFConverter struct{}

var _ Converter = (*MAGICKToHEIFConverter)(nil)

func NewMAGICKToHEIFConverter() *MAGICKToHEIFConverter {
	return &MAGICKToHEIFConverter{}
}

func (c *MAGICKToHEIFConverter) SourceFormat() string {
	return "magick"
}

func (c *MAGICKToHEIFConverter) TargetFormat() string {
	return "heif"
}

func (c *MAGICKToHEIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.HEIF)
}
