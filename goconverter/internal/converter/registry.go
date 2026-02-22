package converter

import "github.com/h2non/bimg"

var converterFactories = []func() Converter{
	func() Converter { return NewAVIFToGIFConverter() },
	func() Converter { return NewAVIFToHEIFConverter() },
	func() Converter { return NewAVIFToJPEGConverter() },
	func() Converter { return NewAVIFToPNGConverter() },
	func() Converter { return NewAVIFToTIFFConverter() },
	func() Converter { return NewAVIFToWEBPConverter() },
	func() Converter { return NewGIFToAVIFConverter() },
	func() Converter { return NewGIFToHEIFConverter() },
	func() Converter { return NewGIFToJPEGConverter() },
	func() Converter { return NewGIFToPNGConverter() },
	func() Converter { return NewGIFToTIFFConverter() },
	func() Converter { return NewGIFToWEBPConverter() },
	func() Converter { return NewHEIFToAVIFConverter() },
	func() Converter { return NewHEIFToGIFConverter() },
	func() Converter { return NewHEIFToJPEGConverter() },
	func() Converter { return NewHEIFToPNGConverter() },
	func() Converter { return NewHEIFToTIFFConverter() },
	func() Converter { return NewHEIFToWEBPConverter() },
	func() Converter { return NewJPEGToAVIFConverter() },
	func() Converter { return NewJPEGToGIFConverter() },
	func() Converter { return NewJPEGToHEIFConverter() },
	func() Converter { return NewJPEGToPNGConverter() },
	func() Converter { return NewJPEGToTIFFConverter() },
	func() Converter { return NewJPEGToWEBPConverter() },
	func() Converter { return NewPNGToAVIFConverter() },
	func() Converter { return NewPNGToGIFConverter() },
	func() Converter { return NewPNGToHEIFConverter() },
	func() Converter { return NewPNGToJPEGConverter() },
	func() Converter { return NewPNGToTIFFConverter() },
	func() Converter { return NewPNGToWEBPConverter() },
	func() Converter { return NewTIFFToAVIFConverter() },
	func() Converter { return NewTIFFToGIFConverter() },
	func() Converter { return NewTIFFToHEIFConverter() },
	func() Converter { return NewTIFFToJPEGConverter() },
	func() Converter { return NewTIFFToPNGConverter() },
	func() Converter { return NewTIFFToWEBPConverter() },
	func() Converter { return NewWEBPToAVIFConverter() },
	func() Converter { return NewWEBPToGIFConverter() },
	func() Converter { return NewWEBPToHEIFConverter() },
	func() Converter { return NewWEBPToJPEGConverter() },
	func() Converter { return NewWEBPToPNGConverter() },
	func() Converter { return NewWEBPToTIFFConverter() },
	func() Converter { return NewMAGICKToAVIFConverter() },
	func() Converter { return NewMAGICKToGIFConverter() },
	func() Converter { return NewMAGICKToHEIFConverter() },
	func() Converter { return NewMAGICKToJPEGConverter() },
	func() Converter { return NewMAGICKToPNGConverter() },
	func() Converter { return NewMAGICKToTIFFConverter() },
	func() Converter { return NewMAGICKToWEBPConverter() },
	func() Converter { return NewPDFToAVIFConverter() },
	func() Converter { return NewPDFToGIFConverter() },
	func() Converter { return NewPDFToHEIFConverter() },
	func() Converter { return NewPDFToJPEGConverter() },
	func() Converter { return NewPDFToPNGConverter() },
	func() Converter { return NewPDFToTIFFConverter() },
	func() Converter { return NewPDFToWEBPConverter() },
	func() Converter { return NewSVGToAVIFConverter() },
	func() Converter { return NewSVGToGIFConverter() },
	func() Converter { return NewSVGToHEIFConverter() },
	func() Converter { return NewSVGToJPEGConverter() },
	func() Converter { return NewSVGToPNGConverter() },
	func() Converter { return NewSVGToTIFFConverter() },
	func() Converter { return NewSVGToWEBPConverter() },
}

var converters = buildConverters()

func buildConverters() []Converter {
	output := make([]Converter, 0, len(converterFactories))
	for _, factory := range converterFactories {
		converter := factory()
		if !bimg.IsTypeNameSupported(converter.SourceFormat()) {
			continue
		}
		if !bimg.IsTypeNameSupportedSave(converter.TargetFormat()) {
			continue
		}
		output = append(output, converter)
	}
	return output
}

func RegisteredConverters() []Converter {
	output := make([]Converter, len(converters))
	copy(output, converters)
	return output
}

func FindConverter(source string, target string) (Converter, bool) {
	for _, c := range converters {
		if c.SourceFormat() == source && c.TargetFormat() == target {
			return c, true
		}
	}

	return nil, false
}

func ConversionTargetsBySource() map[string][]string {
	registered := RegisteredConverters()
	output := make(map[string][]string, len(registered))

	for _, c := range registered {
		source := c.SourceFormat()
		target := c.TargetFormat()
		output[source] = append(output[source], target)
	}

	return output
}
