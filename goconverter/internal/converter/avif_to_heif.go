package converter

import "github.com/h2non/bimg"

type AVIFToHEIFConverter struct{}

var _ Converter = (*AVIFToHEIFConverter)(nil)

func NewAVIFToHEIFConverter() *AVIFToHEIFConverter {
	return &AVIFToHEIFConverter{}
}

func (c *AVIFToHEIFConverter) SourceFormat() string {
	return "avif"
}

func (c *AVIFToHEIFConverter) TargetFormat() string {
	return "heif"
}

func (c *AVIFToHEIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.HEIF)
}
