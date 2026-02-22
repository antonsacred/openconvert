package converter

import "github.com/h2non/bimg"

type PDFToTIFFConverter struct{}

var _ Converter = (*PDFToTIFFConverter)(nil)

func NewPDFToTIFFConverter() *PDFToTIFFConverter {
	return &PDFToTIFFConverter{}
}

func (c *PDFToTIFFConverter) SourceFormat() string {
	return "pdf"
}

func (c *PDFToTIFFConverter) TargetFormat() string {
	return "tiff"
}

func (c *PDFToTIFFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.TIFF)
}
