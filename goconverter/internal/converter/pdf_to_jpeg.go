package converter

import "github.com/h2non/bimg"

type PDFToJPEGConverter struct{}

var _ Converter = (*PDFToJPEGConverter)(nil)

func NewPDFToJPEGConverter() *PDFToJPEGConverter {
	return &PDFToJPEGConverter{}
}

func (c *PDFToJPEGConverter) SourceFormat() string {
	return "pdf"
}

func (c *PDFToJPEGConverter) TargetFormat() string {
	return "jpeg"
}

func (c *PDFToJPEGConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.JPEG)
}
