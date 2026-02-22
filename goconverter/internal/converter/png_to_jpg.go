package converter

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
)

const defaultJPEGQuality = 90

type PNGToJPGConverter struct {
	quality int
}

func NewPNGToJPGConverter() *PNGToJPGConverter {
	return &PNGToJPGConverter{quality: defaultJPEGQuality}
}

func (c *PNGToJPGConverter) SourceFormat() string {
	return "png"
}

func (c *PNGToJPGConverter) TargetFormat() string {
	return "jpg"
}

func (c *PNGToJPGConverter) Convert(input []byte) ([]byte, error) {
	img, err := png.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, fmt.Errorf("decode png: %w", err)
	}

	var output bytes.Buffer
	if err := jpeg.Encode(&output, img, &jpeg.Options{Quality: c.quality}); err != nil {
		return nil, fmt.Errorf("encode jpg: %w", err)
	}

	return output.Bytes(), nil
}
