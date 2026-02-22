package converter

import (
	"fmt"

	"github.com/h2non/bimg"
)

func convertWithBIMG(input []byte, sourceFormat string, targetFormat string, targetType bimg.ImageType) ([]byte, error) {
	output, err := bimg.NewImage(input).Convert(targetType)
	if err != nil {
		return nil, fmt.Errorf("convert %s to %s: %w", sourceFormat, targetFormat, err)
	}

	return output, nil
}
