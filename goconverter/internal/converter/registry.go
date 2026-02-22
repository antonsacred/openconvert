package converter

var converters = []Converter{
	NewPNGToJPGConverter(),
}

func RegisteredConverters() []Converter {
	output := make([]Converter, len(converters))
	copy(output, converters)
	return output
}

func PossibleConversions() [][2]string {
	registered := RegisteredConverters()
	output := make([][2]string, 0, len(registered))
	for _, c := range registered {
		output = append(output, [2]string{c.SourceFormat(), c.TargetFormat()})
	}
	return output
}
