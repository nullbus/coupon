package coupon

import (
	"errors"
	"regexp"
	"strings"
)

// Validator figures out the given code is exact.
type Validator struct {
	NumParts int
}

func (v *Validator) Validate(code string) (string, error) {
	numParts := v.NumParts
	if numParts == 0 {
		numParts = defaultNumParts
	}

	code = strings.ToUpper(code)
	code = regexp.MustCompile("[^0-9A-Z]+").ReplaceAllString(code, "")
	code = strings.Map(func(r rune) rune {
		if idx := strings.IndexRune("OIZS", r); idx >= 0 {
			return rune("0125"[idx])
		}
		return r
	}, code)

	parts := regexp.MustCompile("([0-9A-Z]){4}").FindAllString(code, -1)
	if len(parts) != numParts {
		return "", errors.New("number of parts not match")
	}

	for i, part := range parts {
		str, check := []byte(part[:3]), part[3]
		if checkDigitAlg1(str, i+1) != check {
			return "", errors.New("invalid checksum")
		}
	}

	return strings.Join(parts, "-"), nil
}
