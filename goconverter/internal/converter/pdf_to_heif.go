package converter

import "github.com/h2non/bimg"

type PDFToHEIFConverter struct{}

var _ Converter = (*PDFToHEIFConverter)(nil)

func NewPDFToHEIFConverter() *PDFToHEIFConverter {
	return &PDFToHEIFConverter{}
}

func (c *PDFToHEIFConverter) SourceFormat() string {
	return "pdf"
}

func (c *PDFToHEIFConverter) TargetFormat() string {
	return "heif"
}

func (c *PDFToHEIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.HEIF)
}
