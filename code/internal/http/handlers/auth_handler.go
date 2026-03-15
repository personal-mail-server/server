package handlers

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"personal-mail-server/internal/auth"
)

type AuthHandler struct {
	service *auth.Service
}

func NewAuthHandler(service *auth.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req auth.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, auth.ErrorResponse{
			Code:    auth.CodeInvalidRequestBody,
			Message: "입력값 형식이 올바르지 않습니다.",
		})
	}

	resp, appErr := h.service.Login(c.Request().Context(), req)
	if appErr != nil {
		return c.JSON(appErr.Status, auth.ErrorResponse{
			Code:    appErr.Code,
			Message: appErr.Message,
		})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Logout(c echo.Context) error {
	rawToken, ok := extractBearerToken(c.Request().Header.Get(echo.HeaderAuthorization))
	if !ok {
		appErr := auth.NewInvalidAccessToken()
		return c.JSON(appErr.Status, auth.ErrorResponse{Code: appErr.Code, Message: appErr.Message})
	}

	appErr := h.service.Logout(c.Request().Context(), rawToken)
	if appErr != nil {
		return c.JSON(appErr.Status, auth.ErrorResponse{Code: appErr.Code, Message: appErr.Message})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *AuthHandler) Reissue(c echo.Context) error {
	var req auth.ReissueRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, auth.ErrorResponse{
			Code:    auth.CodeInvalidRequestBody,
			Message: "입력값 형식이 올바르지 않습니다.",
		})
	}

	resp, appErr := h.service.Reissue(c.Request().Context(), req)
	if appErr != nil {
		return c.JSON(appErr.Status, auth.ErrorResponse{Code: appErr.Code, Message: appErr.Message})
	}

	return c.JSON(http.StatusOK, resp)
}

func extractBearerToken(header string) (string, bool) {
	prefix := auth.TokenTypeBearer + " "
	if !strings.HasPrefix(header, prefix) {
		return "", false
	}
	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if token == "" {
		return "", false
	}
	return token, true
}
