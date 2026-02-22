package converter

import "github.com/h2non/bimg"

type WEBPToJPGConverter struct{}

var _ Converter = (*WEBPToJPGConverter)(nil)

func NewWEBPToJPGConverter() *WEBPToJPGConverter {
	return &WEBPToJPGConverter{}
}

func (c *WEBPToJPGConverter) SourceFormat() string {
	return "webp"
}

func (c *WEBPToJPGConverter) TargetFormat() string {
	return "jpg"
}

func (c *WEBPToJPGConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.JPEG)
}
