package nspv

import (
	"context"

	"github.com/agnivade/levenshtein"
)

const (
	DefaultMinLength            = 8  // Default minimum length of the password.
	DefaultMaxLength            = 64 // Default maximum length of the password.
	DefaultHibpThreshold        = 0  // Default threshold of HIBP count for the password.
	DefaultLevenshteinThreshold = 2  // Default threshold of the Levenshtein distance to the dictionary for the password.
)

// Validator
type Validator struct {
	// Dictionary of the request context.
	dict []string

	// minimum length of the password.
	minLength int

	// maximum length of the password.
	maxLength int

	// threshold of HIBP count for the password.
	hibpThreshold int

	// threshold of the Levenshtein distance to the dictionary for the password.
	levenshteinThreshold int

	// the flag for ignore hibp error.
	ignoreHibpError bool

	hc *hibpClient
}

// NewValidator returns Validator.
func NewValidator() *Validator {
	v := Validator{}
	v.minLength = DefaultMinLength
	v.maxLength = DefaultMaxLength
	v.hibpThreshold = DefaultHibpThreshold
	v.levenshteinThreshold = DefaultLevenshteinThreshold

	v.hc = newHibpClient()
	v.hc.ctx = context.TODO()
	return &v
}

// SetDict set a dictionary composed of request context.
func (v *Validator) SetDict(dict []string) {
	v.dict = dict
}

// SetMinLength set minimum length of the password.
func (v *Validator) SetMinLength(length int) {
	v.minLength = length
}

// SetMaxLength set maximum length of the password.
func (v *Validator) SetMaxLength(length int) {
	v.maxLength = length
}

// SetHibpThreshold set threshold of HIBP count for the password.
func (v *Validator) SetHibpThreshold(threshold int) {
	v.hibpThreshold = threshold
}

// SetLevenshteinThreshold set threshold of the Levenshtein distance to the dictionary for the password.
func (v *Validator) SetLevenshteinThreshold(threshold int) {
	v.levenshteinThreshold = threshold
}

// SetHibpClientContext set context.Context for the request to HIBP.
func (v *Validator) SetHibpClientContext(ctx context.Context) {
	v.hc.ctx = ctx
}

// SetIgnoreHibpError set the flag for ignore hibp error. (not recommended)
func (v *Validator) SetIgnoreHibpError(flag bool) {
	v.ignoreHibpError = flag
}

// Validate validates the password.
func (v *Validator) Validate(password string) (result Result, err error) {
	result, err = v.checkMinLength(password)
	if result != Ok || err != nil {
		return
	}

	result, err = v.checkMaxLength(password)
	if result != Ok || err != nil {
		return
	}

	result, err = v.checkDict(password)
	if result != Ok || err != nil {
		return
	}

	result, err = v.checkHibp(password)
	if result != Ok || err != nil {
		return
	}

	return Ok, nil
}

func (v *Validator) checkMinLength(password string) (Result, error) {
	if len(password) < v.minLength {
		return ViolateMinLengthCheck, nil
	}
	return Ok, nil
}

func (v *Validator) checkMaxLength(password string) (Result, error) {
	if v.maxLength < len(password) {
		return ViolateMaxLengthCheck, nil
	}
	return Ok, nil
}

func (v *Validator) checkDict(password string) (Result, error) {
	for _, word := range v.dict {
		if levenshtein.ComputeDistance(password, word) < v.levenshteinThreshold {
			return ViolateDictCheck, nil
		}
	}
	return Ok, nil
}

func (v *Validator) checkHibp(password string) (Result, error) {
	count, err := v.hc.pwnedCount(password)
	if err != nil {
		if v.ignoreHibpError {
			return Ok, nil
		}

		return Error, err
	}

	if count > v.hibpThreshold {
		return ViolateHibpCheck, nil
	}
	return Ok, nil
}
