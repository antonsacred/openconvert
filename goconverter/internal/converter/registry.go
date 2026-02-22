package converter

var converters = []Converter{
	NewPNGToJPGConverter(),
	NewPNGToWEBPConverter(),
	NewJPGToPNGConverter(),
	NewJPGToWEBPConverter(),
	NewWEBPToJPGConverter(),
	NewWEBPToPNGConverter(),
}

func RegisteredConverters() []Converter {
	output := make([]Converter, len(converters))
	copy(output, converters)
	return output
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
