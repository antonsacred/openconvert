package converter

import "github.com/h2non/bimg"

type TIFFToJPEGConverter struct{}

var _ Converter = (*TIFFToJPEGConverter)(nil)

func NewTIFFToJPEGConverter() *TIFFToJPEGConverter {
	return &TIFFToJPEGConverter{}
}

func (c *TIFFToJPEGConverter) SourceFormat() string {
	return "tiff"
}

func (c *TIFFToJPEGConverter) TargetFormat() string {
	return "jpeg"
}

func (c *TIFFToJPEGConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.JPEG)
}
