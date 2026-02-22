package converter

var supportedConversions = [][2]string{
	{"jpg", "png"},
	{"png", "jpg"},
	{"jpg", "webp"},
	{"webp", "jpg"},
	{"png", "webp"},
	{"webp", "png"},
	{"txt", "pdf"},
	{"md", "html"},
	{"docx", "pdf"},
}

func PossibleConversions() [][2]string {
	output := make([][2]string, len(supportedConversions))
	copy(output, supportedConversions)
	return output
}
