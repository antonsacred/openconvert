package converter

import "github.com/h2non/bimg"

type AVIFToWEBPConverter struct{}

var _ Converter = (*AVIFToWEBPConverter)(nil)

func NewAVIFToWEBPConverter() *AVIFToWEBPConverter {
	return &AVIFToWEBPConverter{}
}

func (c *AVIFToWEBPConverter) SourceFormat() string {
	return "avif"
}

func (c *AVIFToWEBPConverter) TargetFormat() string {
	return "webp"
}

func (c *AVIFToWEBPConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.WEBP)
}
