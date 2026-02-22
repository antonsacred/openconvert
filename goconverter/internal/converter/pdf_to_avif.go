package converter

import "github.com/h2non/bimg"

type PDFToAVIFConverter struct{}

var _ Converter = (*PDFToAVIFConverter)(nil)

func NewPDFToAVIFConverter() *PDFToAVIFConverter {
	return &PDFToAVIFConverter{}
}

func (c *PDFToAVIFConverter) SourceFormat() string {
	return "pdf"
}

func (c *PDFToAVIFConverter) TargetFormat() string {
	return "avif"
}

func (c *PDFToAVIFConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.AVIF)
}
