package auth

import "net/http"

const (
	CodeInvalidRequestBody  = "INVALID_REQUEST_BODY"
	CodeInvalidLoginID      = "INVALID_LOGIN_ID_FORMAT"
	CodeInvalidPassword     = "INVALID_PASSWORD_FORMAT"
	CodeMissingRequired     = "MISSING_REQUIRED_FIELD"
	CodeInvalidCredentials  = "INVALID_CREDENTIALS"
	CodeInvalidAccessToken  = "INVALID_ACCESS_TOKEN"
	CodeInvalidRefreshToken = "INVALID_REFRESH_TOKEN"
	CodeAccountLocked       = "ACCOUNT_LOCKED"
	CodeInternalServerError = "INTERNAL_SERVER_ERROR"
)

type AppError struct {
	Status  int
	Code    string
	Message string
}

func (e *AppError) Error() string {
	return e.Code
}

func NewBadRequest(code string) *AppError {
	return &AppError{Status: http.StatusBadRequest, Code: code, Message: "입력값 형식이 올바르지 않습니다."}
}

func NewUnauthorized() *AppError {
	return &AppError{Status: http.StatusUnauthorized, Code: CodeInvalidCredentials, Message: "로그인 ID 또는 비밀번호가 올바르지 않습니다."}
}

func NewInvalidAccessToken() *AppError {
	return &AppError{Status: http.StatusUnauthorized, Code: CodeInvalidAccessToken, Message: "인증 정보가 유효하지 않습니다. 다시 로그인해 주세요."}
}

func NewInvalidRefreshToken() *AppError {
	return &AppError{Status: http.StatusUnauthorized, Code: CodeInvalidRefreshToken, Message: "인증 정보가 유효하지 않습니다. 다시 로그인해 주세요."}
}

func NewLocked() *AppError {
	return &AppError{Status: http.StatusLocked, Code: CodeAccountLocked, Message: "계정이 일시적으로 잠겼습니다. 잠시 후 다시 시도해 주세요."}
}

func NewInternalServerError() *AppError {
	return &AppError{Status: http.StatusInternalServerError, Code: CodeInternalServerError, Message: "요청을 처리할 수 없습니다. 잠시 후 다시 시도해 주세요."}
}
