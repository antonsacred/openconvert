package converter

type Converter interface {
	SourceFormat() string
	TargetFormat() string
	Convert(input []byte) ([]byte, error)
}
