package nspv

// Result of the validation.
type Result int

const (
	_                     Result = iota
	Ok                           // Validation OK
	ViolateMinLengthCheck        // Violate minimum length check.
	ViolateMaxLengthCheck        // Violate maximum length check.
	ViolateDictCheck             // Violate dictionary check.
	ViolateHibpCheck             // Violate HIBP check.
	Error                        // Validation Error
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
	case ViolateHibpCheck:
		return "ViolateHibpCheck"
	case Error:
		return "Error"
	default:
		return "Unknown"
	}
}
