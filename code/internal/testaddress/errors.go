package testaddress

import "errors"

var (
	ErrTestMailAddressNotFound = errors.New("test mail address not found")
	ErrDuplicateEmail          = errors.New("duplicate test mail address email")
)
