package converter

import "github.com/h2non/bimg"

type TIFFToPNGConverter struct{}

var _ Converter = (*TIFFToPNGConverter)(nil)

func NewTIFFToPNGConverter() *TIFFToPNGConverter {
	return &TIFFToPNGConverter{}
}

func (c *TIFFToPNGConverter) SourceFormat() string {
	return "tiff"
}

func (c *TIFFToPNGConverter) TargetFormat() string {
	return "png"
}

func (c *TIFFToPNGConverter) Convert(input []byte) ([]byte, error) {
	return convertWithBIMG(input, c.SourceFormat(), c.TargetFormat(), bimg.PNG)
}
