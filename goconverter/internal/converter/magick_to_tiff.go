package converter

import "github.com/h2non/bimg"

type MAGICKToTIFFConverter struct{}

var _ Converter = (*MAGICKToTIFFConverter)(nil)

func NewMAGICKToTIFFConverter() *MAGICKToTIFFConverter {
	return &MAGICKToTIFFConverter{}
}

func (c *MAGICKToTIFFConverter) SourceFormat() string {
	return "magick"
}

func (c *MAGICKToTIFFConverter) TargetFormat() string {
	return "tiff"
}

func (c *MAGICKToTIFFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.TIFF)
}
