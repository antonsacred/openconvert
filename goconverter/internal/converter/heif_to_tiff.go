package converter

import "github.com/h2non/bimg"

type HEIFToTIFFConverter struct{}

var _ Converter = (*HEIFToTIFFConverter)(nil)

func NewHEIFToTIFFConverter() *HEIFToTIFFConverter {
	return &HEIFToTIFFConverter{}
}

func (c *HEIFToTIFFConverter) SourceFormat() string {
	return "heif"
}

func (c *HEIFToTIFFConverter) TargetFormat() string {
	return "tiff"
}

func (c *HEIFToTIFFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.TIFF)
}
