package converter

import "github.com/h2non/bimg"

type MAGICKToWEBPConverter struct{}

var _ Converter = (*MAGICKToWEBPConverter)(nil)

func NewMAGICKToWEBPConverter() *MAGICKToWEBPConverter {
	return &MAGICKToWEBPConverter{}
}

func (c *MAGICKToWEBPConverter) SourceFormat() string {
	return "magick"
}

func (c *MAGICKToWEBPConverter) TargetFormat() string {
	return "webp"
}

func (c *MAGICKToWEBPConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.WEBP)
}
