package converter

import "github.com/h2non/bimg"

type HEIFToPNGConverter struct{}

var _ Converter = (*HEIFToPNGConverter)(nil)

func NewHEIFToPNGConverter() *HEIFToPNGConverter {
	return &HEIFToPNGConverter{}
}

func (c *HEIFToPNGConverter) SourceFormat() string {
	return "heif"
}

func (c *HEIFToPNGConverter) TargetFormat() string {
	return "png"
}

func (c *HEIFToPNGConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.PNG)
}
