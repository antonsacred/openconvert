package converter

import "github.com/h2non/bimg"

type JPGToPNGConverter struct{}

var _ Converter = (*JPGToPNGConverter)(nil)

func NewJPGToPNGConverter() *JPGToPNGConverter {
	return &JPGToPNGConverter{}
}

func (c *JPGToPNGConverter) SourceFormat() string {
	return "jpg"
}

func (c *JPGToPNGConverter) TargetFormat() string {
	return "png"
}

func (c *JPGToPNGConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.PNG)
}
