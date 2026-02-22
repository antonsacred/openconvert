package converter

import "github.com/h2non/bimg"

type WEBPToGIFConverter struct{}

var _ Converter = (*WEBPToGIFConverter)(nil)

func NewWEBPToGIFConverter() *WEBPToGIFConverter {
	return &WEBPToGIFConverter{}
}

func (c *WEBPToGIFConverter) SourceFormat() string {
	return "webp"
}

func (c *WEBPToGIFConverter) TargetFormat() string {
	return "gif"
}

func (c *WEBPToGIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.GIF)
}
