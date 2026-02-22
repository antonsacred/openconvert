package converter

import "github.com/h2non/bimg"

type WEBPToJPEGConverter struct{}

var _ Converter = (*WEBPToJPEGConverter)(nil)

func NewWEBPToJPEGConverter() *WEBPToJPEGConverter {
	return &WEBPToJPEGConverter{}
}

func (c *WEBPToJPEGConverter) SourceFormat() string {
	return "webp"
}

func (c *WEBPToJPEGConverter) TargetFormat() string {
	return "jpeg"
}

func (c *WEBPToJPEGConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.JPEG)
}
