package service

import (
	"regexp"

	"gopkg.in/go-playground/validator.v9"
)

type (
	// CustomValidator - CustomValidator
	CustomValidator struct {
		Validator *validator.Validate
	}
)

// Validate - Validate HTTP Data Input
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

var isStringAlphabetic = regexp.MustCompile(`^[a-zA-Z0-9_]*$`).MatchString
var isPhoneNumber = regexp.MustCompile(`^[0-9_]*$`).MatchString
