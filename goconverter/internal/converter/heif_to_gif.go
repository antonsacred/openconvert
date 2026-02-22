package converter

import "github.com/h2non/bimg"

type HEIFToGIFConverter struct{}

var _ Converter = (*HEIFToGIFConverter)(nil)

func NewHEIFToGIFConverter() *HEIFToGIFConverter {
	return &HEIFToGIFConverter{}
}

func (c *HEIFToGIFConverter) SourceFormat() string {
	return "heif"
}

func (c *HEIFToGIFConverter) TargetFormat() string {
	return "gif"
}

func (c *HEIFToGIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.GIF)
}
