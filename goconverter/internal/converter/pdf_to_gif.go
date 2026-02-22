package converter

import "github.com/h2non/bimg"

type PDFToGIFConverter struct{}

var _ Converter = (*PDFToGIFConverter)(nil)

func NewPDFToGIFConverter() *PDFToGIFConverter {
	return &PDFToGIFConverter{}
}

func (c *PDFToGIFConverter) SourceFormat() string {
	return "pdf"
}

func (c *PDFToGIFConverter) TargetFormat() string {
	return "gif"
}

func (c *PDFToGIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.GIF)
}
