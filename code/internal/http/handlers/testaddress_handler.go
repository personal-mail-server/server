package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"personal-mail-server/internal/auth"
	"personal-mail-server/internal/testaddress"
)

type TestAddressHandler struct {
	service *testaddress.Service
}

func NewTestAddressHandler(service *testaddress.Service) *TestAddressHandler {
	return &TestAddressHandler{service: service}
}

func (h *TestAddressHandler) GenerateCandidate(c echo.Context) error {
	rawToken, ok := extractBearerToken(c.Request().Header.Get(echo.HeaderAuthorization))
	if !ok {
		appErr := auth.NewInvalidAccessToken()
		return c.JSON(appErr.Status, auth.ErrorResponse{Code: appErr.Code, Message: appErr.Message})
	}

	resp, appErr := h.service.GenerateCandidate(c.Request().Context(), rawToken)
	if appErr != nil {
		return c.JSON(appErr.Status, auth.ErrorResponse{Code: appErr.Code, Message: appErr.Message})
	}

	return c.JSON(http.StatusOK, resp)
}
