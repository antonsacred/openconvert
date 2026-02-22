package converter

import "github.com/h2non/bimg"

type MAGICKToAVIFConverter struct{}

var _ Converter = (*MAGICKToAVIFConverter)(nil)

func NewMAGICKToAVIFConverter() *MAGICKToAVIFConverter {
	return &MAGICKToAVIFConverter{}
}

func (c *MAGICKToAVIFConverter) SourceFormat() string {
	return "magick"
}

func (c *MAGICKToAVIFConverter) TargetFormat() string {
	return "avif"
}

func (c *MAGICKToAVIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.AVIF)
}
