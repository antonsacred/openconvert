package converter

import "github.com/h2non/bimg"

type GIFToTIFFConverter struct{}

var _ Converter = (*GIFToTIFFConverter)(nil)

func NewGIFToTIFFConverter() *GIFToTIFFConverter {
	return &GIFToTIFFConverter{}
}

func (c *GIFToTIFFConverter) SourceFormat() string {
	return "gif"
}

func (c *GIFToTIFFConverter) TargetFormat() string {
	return "tiff"
}

func (c *GIFToTIFFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.TIFF)
}
