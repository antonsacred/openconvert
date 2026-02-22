package converter

import "github.com/h2non/bimg"

type WEBPToHEIFConverter struct{}

var _ Converter = (*WEBPToHEIFConverter)(nil)

func NewWEBPToHEIFConverter() *WEBPToHEIFConverter {
	return &WEBPToHEIFConverter{}
}

func (c *WEBPToHEIFConverter) SourceFormat() string {
	return "webp"
}

func (c *WEBPToHEIFConverter) TargetFormat() string {
	return "heif"
}

func (c *WEBPToHEIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.HEIF)
}
