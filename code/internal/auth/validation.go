package auth

import (
	"regexp"
	"unicode"
	"unicode/utf8"
)

var loginIDPattern = regexp.MustCompile(`^[a-z0-9-]{4,32}$`)

func ValidateLoginRequest(req LoginRequest) *AppError {
	if req.LoginID == "" || req.Password == "" {
		return NewBadRequest(CodeMissingRequired)
	}

	if !loginIDPattern.MatchString(req.LoginID) {
		return NewBadRequest(CodeInvalidLoginID)
	}

	passwordLength := utf8.RuneCountInString(req.Password)
	if passwordLength < 8 || passwordLength > 64 {
		return NewBadRequest(CodeInvalidPassword)
	}

	hasLetter := false
	hasDigit := false
	for _, r := range req.Password {
		if unicode.IsSpace(r) {
			return NewBadRequest(CodeInvalidPassword)
		}
		if unicode.IsLetter(r) {
			hasLetter = true
		}
		if unicode.IsDigit(r) {
			hasDigit = true
		}
	}

	if !hasLetter || !hasDigit {
		return NewBadRequest(CodeInvalidPassword)
	}

	return nil
}
