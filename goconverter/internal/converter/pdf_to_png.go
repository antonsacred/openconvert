package converter

import "github.com/h2non/bimg"

type PDFToPNGConverter struct{}

var _ Converter = (*PDFToPNGConverter)(nil)

func NewPDFToPNGConverter() *PDFToPNGConverter {
	return &PDFToPNGConverter{}
}

func (c *PDFToPNGConverter) SourceFormat() string {
	return "pdf"
}

func (c *PDFToPNGConverter) TargetFormat() string {
	return "png"
}

func (c *PDFToPNGConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.PNG)
}
