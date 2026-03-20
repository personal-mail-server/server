package testaddress

import (
	"net/mail"
	"strings"

	"personal-mail-server/internal/auth"
)

func ValidateCreateRequest(req CreateRequest) *auth.AppError {
	if strings.TrimSpace(req.Email) == "" {
		return auth.NewBadRequest(auth.CodeMissingRequired)
	}

	parsed, err := mail.ParseAddress(req.Email)
	if err != nil || parsed.Address != req.Email {
		return auth.NewBadRequest(auth.CodeInvalidEmail)
	}

	return nil
}
