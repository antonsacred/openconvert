package converter

import "github.com/h2non/bimg"

type MAGICKToPNGConverter struct{}

var _ Converter = (*MAGICKToPNGConverter)(nil)

func NewMAGICKToPNGConverter() *MAGICKToPNGConverter {
	return &MAGICKToPNGConverter{}
}

func (c *MAGICKToPNGConverter) SourceFormat() string {
	return "magick"
}

func (c *MAGICKToPNGConverter) TargetFormat() string {
	return "png"
}

func (c *MAGICKToPNGConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.PNG)
}
