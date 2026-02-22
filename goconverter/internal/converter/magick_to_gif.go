package converter

import "github.com/h2non/bimg"

type MAGICKToGIFConverter struct{}

var _ Converter = (*MAGICKToGIFConverter)(nil)

func NewMAGICKToGIFConverter() *MAGICKToGIFConverter {
	return &MAGICKToGIFConverter{}
}

func (c *MAGICKToGIFConverter) SourceFormat() string {
	return "magick"
}

func (c *MAGICKToGIFConverter) TargetFormat() string {
	return "gif"
}

func (c *MAGICKToGIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.GIF)
}
