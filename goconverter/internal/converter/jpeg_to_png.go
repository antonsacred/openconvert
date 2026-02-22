package converter

import "github.com/h2non/bimg"

type JPEGToPNGConverter struct{}

var _ Converter = (*JPEGToPNGConverter)(nil)

func NewJPEGToPNGConverter() *JPEGToPNGConverter {
	return &JPEGToPNGConverter{}
}

func (c *JPEGToPNGConverter) SourceFormat() string {
	return "jpeg"
}

func (c *JPEGToPNGConverter) TargetFormat() string {
	return "png"
}

func (c *JPEGToPNGConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.PNG)
}
