package converter

import "github.com/h2non/bimg"

type AVIFToGIFConverter struct{}

var _ Converter = (*AVIFToGIFConverter)(nil)

func NewAVIFToGIFConverter() *AVIFToGIFConverter {
	return &AVIFToGIFConverter{}
}

func (c *AVIFToGIFConverter) SourceFormat() string {
	return "avif"
}

func (c *AVIFToGIFConverter) TargetFormat() string {
	return "gif"
}

func (c *AVIFToGIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.GIF)
}
