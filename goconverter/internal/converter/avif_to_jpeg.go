package converter

import "github.com/h2non/bimg"

type AVIFToJPEGConverter struct{}

var _ Converter = (*AVIFToJPEGConverter)(nil)

func NewAVIFToJPEGConverter() *AVIFToJPEGConverter {
	return &AVIFToJPEGConverter{}
}

func (c *AVIFToJPEGConverter) SourceFormat() string {
	return "avif"
}

func (c *AVIFToJPEGConverter) TargetFormat() string {
	return "jpeg"
}

func (c *AVIFToJPEGConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.JPEG)
}
