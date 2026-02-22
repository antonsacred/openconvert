package converter

import "github.com/h2non/bimg"

type WEBPToPNGConverter struct{}

var _ Converter = (*WEBPToPNGConverter)(nil)

func NewWEBPToPNGConverter() *WEBPToPNGConverter {
	return &WEBPToPNGConverter{}
}

func (c *WEBPToPNGConverter) SourceFormat() string {
	return "webp"
}

func (c *WEBPToPNGConverter) TargetFormat() string {
	return "png"
}

func (c *WEBPToPNGConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.PNG)
}
