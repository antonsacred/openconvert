package converter

import "github.com/h2non/bimg"

type MAGICKToJPEGConverter struct{}

var _ Converter = (*MAGICKToJPEGConverter)(nil)

func NewMAGICKToJPEGConverter() *MAGICKToJPEGConverter {
	return &MAGICKToJPEGConverter{}
}

func (c *MAGICKToJPEGConverter) SourceFormat() string {
	return "magick"
}

func (c *MAGICKToJPEGConverter) TargetFormat() string {
	return "jpeg"
}

func (c *MAGICKToJPEGConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.JPEG)
}
