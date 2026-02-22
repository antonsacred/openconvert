package converter

import "github.com/h2non/bimg"

type JPEGToTIFFConverter struct{}

var _ Converter = (*JPEGToTIFFConverter)(nil)

func NewJPEGToTIFFConverter() *JPEGToTIFFConverter {
	return &JPEGToTIFFConverter{}
}

func (c *JPEGToTIFFConverter) SourceFormat() string {
	return "jpeg"
}

func (c *JPEGToTIFFConverter) TargetFormat() string {
	return "tiff"
}

func (c *JPEGToTIFFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.TIFF)
}
