package converter

import "github.com/h2non/bimg"

type AVIFToTIFFConverter struct{}

var _ Converter = (*AVIFToTIFFConverter)(nil)

func NewAVIFToTIFFConverter() *AVIFToTIFFConverter {
	return &AVIFToTIFFConverter{}
}

func (c *AVIFToTIFFConverter) SourceFormat() string {
	return "avif"
}

func (c *AVIFToTIFFConverter) TargetFormat() string {
	return "tiff"
}

func (c *AVIFToTIFFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.TIFF)
}
