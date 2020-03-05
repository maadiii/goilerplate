package app

import (
	"errors"
	"unicode"

	"github.com/google/uuid"
)

func CheckUUID(value interface{}) error {
	id := value.(uuid.UUID)
	if id == uuid.Nil {
		return errors.New("invalid uuid")
	}
	return nil
}

func IsPasswordValid(s string) bool {
	var (
		hasMinLen = false
		hasUpper  = false
		hasLower  = false
		hasNumber = false
	)
	if len(s) >= 6 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber
}
