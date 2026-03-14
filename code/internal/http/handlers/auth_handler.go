package handlers

import (
	"net/http"

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
