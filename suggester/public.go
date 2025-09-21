package suggester

// ConvertJSON performs pure passthrough JSON transformation.
// Takes arbitrary JSON input, computes positions for LM fields, and passes everything else through.
func ConvertJSON(inputJSON []byte, format string) ([]byte, error) {
	return PassthroughConvert(inputJSON, format)
}