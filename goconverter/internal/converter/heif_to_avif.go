package converter

import "github.com/h2non/bimg"

type HEIFToAVIFConverter struct{}

var _ Converter = (*HEIFToAVIFConverter)(nil)

func NewHEIFToAVIFConverter() *HEIFToAVIFConverter {
	return &HEIFToAVIFConverter{}
}

func (c *HEIFToAVIFConverter) SourceFormat() string {
	return "heif"
}

func (c *HEIFToAVIFConverter) TargetFormat() string {
	return "avif"
}

func (c *HEIFToAVIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.AVIF)
}
