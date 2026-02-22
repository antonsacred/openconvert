package converter

import "github.com/h2non/bimg"

type WEBPToAVIFConverter struct{}

var _ Converter = (*WEBPToAVIFConverter)(nil)

func NewWEBPToAVIFConverter() *WEBPToAVIFConverter {
	return &WEBPToAVIFConverter{}
}

func (c *WEBPToAVIFConverter) SourceFormat() string {
	return "webp"
}

func (c *WEBPToAVIFConverter) TargetFormat() string {
	return "avif"
}

func (c *WEBPToAVIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.AVIF)
}
