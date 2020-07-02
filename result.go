package nspv

type Result int

const (
	Ok Result = iota
	ViolateMinLengthCheck
	ViolateMaxLengthCheck
	ViolateDictCheck
	ViolateBibpCheck
	Error
)

func (r Result) String() string {
	switch r {
	case Ok:
		return "Ok"
	case ViolateMinLengthCheck:
		return "ViolateMinLengthCheck"
	case ViolateMaxLengthCheck:
		return "ViolateMaxLengthCheck"
	case ViolateDictCheck:
		return "ViolateDictCheck"
	case ViolateBibpCheck:
		return "ViolateBibpCheck"
	case Error:
		return "Error"
	default:
		return "Unknown"
	}
}
