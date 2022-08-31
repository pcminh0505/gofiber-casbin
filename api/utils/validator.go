package utils

import (
	valid "github.com/asaskevich/govalidator"
)

// IsEmpty checks if a string is empty
func IsEmpty(str string) bool {
	return valid.HasWhitespaceOnly(str) || str == ""
}
