package converter

import "github.com/h2non/bimg"

type HEIFToWEBPConverter struct{}

var _ Converter = (*HEIFToWEBPConverter)(nil)

func NewHEIFToWEBPConverter() *HEIFToWEBPConverter {
	return &HEIFToWEBPConverter{}
}

func (c *HEIFToWEBPConverter) SourceFormat() string {
	return "heif"
}

func (c *HEIFToWEBPConverter) TargetFormat() string {
	return "webp"
}

func (c *HEIFToWEBPConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.WEBP)
}
