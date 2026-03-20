package testaddress

import (
	"net/mail"
	"strings"

	"personal-mail-server/internal/auth"
)

func ValidateCreateRequest(req CreateRequest) *auth.AppError {
	return validateEmailField(req.Email)
}

func ValidateUpdateRequest(req UpdateRequest) *auth.AppError {
	return validateEmailField(req.Email)
}

func validateEmailField(email string) *auth.AppError {
	if strings.TrimSpace(email) == "" {
		return auth.NewBadRequest(auth.CodeMissingRequired)
	}

	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email {
		return auth.NewBadRequest(auth.CodeInvalidEmail)
	}

	return nil
}
