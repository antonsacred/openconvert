package converter

import "github.com/h2non/bimg"

type WEBPToTIFFConverter struct{}

var _ Converter = (*WEBPToTIFFConverter)(nil)

func NewWEBPToTIFFConverter() *WEBPToTIFFConverter {
	return &WEBPToTIFFConverter{}
}

func (c *WEBPToTIFFConverter) SourceFormat() string {
	return "webp"
}

func (c *WEBPToTIFFConverter) TargetFormat() string {
	return "tiff"
}

func (c *WEBPToTIFFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.TIFF)
}
