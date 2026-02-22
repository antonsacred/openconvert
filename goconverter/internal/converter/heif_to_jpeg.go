package converter

import "github.com/h2non/bimg"

type HEIFToJPEGConverter struct{}

var _ Converter = (*HEIFToJPEGConverter)(nil)

func NewHEIFToJPEGConverter() *HEIFToJPEGConverter {
	return &HEIFToJPEGConverter{}
}

func (c *HEIFToJPEGConverter) SourceFormat() string {
	return "heif"
}

func (c *HEIFToJPEGConverter) TargetFormat() string {
	return "jpeg"
}

func (c *HEIFToJPEGConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.JPEG)
}
