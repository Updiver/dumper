package dumper

func PositiveBoolRef() *bool {
	var posBool = true
	return &posBool
}

func NegativeBoolRef() *bool {
	var negBool = false
	return &negBool
}
