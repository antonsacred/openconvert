package converter

import "github.com/h2non/bimg"

type PDFToWEBPConverter struct{}

var _ Converter = (*PDFToWEBPConverter)(nil)

func NewPDFToWEBPConverter() *PDFToWEBPConverter {
	return &PDFToWEBPConverter{}
}

func (c *PDFToWEBPConverter) SourceFormat() string {
	return "pdf"
}

func (c *PDFToWEBPConverter) TargetFormat() string {
	return "webp"
}

func (c *PDFToWEBPConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.WEBP)
}
