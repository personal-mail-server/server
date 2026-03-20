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

func (h *TestAddressHandler) Create(c echo.Context) error {
	rawToken, ok := extractBearerToken(c.Request().Header.Get(echo.HeaderAuthorization))
	if !ok {
		appErr := auth.NewInvalidAccessToken()
		return c.JSON(appErr.Status, auth.ErrorResponse{Code: appErr.Code, Message: appErr.Message})
	}

	var req testaddress.CreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, auth.ErrorResponse{
			Code:    auth.CodeInvalidRequestBody,
			Message: "입력값 형식이 올바르지 않습니다.",
		})
	}

	resp, appErr := h.service.Create(c.Request().Context(), rawToken, req)
	if appErr != nil {
		return c.JSON(appErr.Status, auth.ErrorResponse{Code: appErr.Code, Message: appErr.Message})
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *TestAddressHandler) List(c echo.Context) error {
	rawToken, ok := extractBearerToken(c.Request().Header.Get(echo.HeaderAuthorization))
	if !ok {
		appErr := auth.NewInvalidAccessToken()
		return c.JSON(appErr.Status, auth.ErrorResponse{Code: appErr.Code, Message: appErr.Message})
	}

	resp, appErr := h.service.List(c.Request().Context(), rawToken)
	if appErr != nil {
		return c.JSON(appErr.Status, auth.ErrorResponse{Code: appErr.Code, Message: appErr.Message})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *TestAddressHandler) GetByID(c echo.Context) error {
	rawToken, ok := extractBearerToken(c.Request().Header.Get(echo.HeaderAuthorization))
	if !ok {
		appErr := auth.NewInvalidAccessToken()
		return c.JSON(appErr.Status, auth.ErrorResponse{Code: appErr.Code, Message: appErr.Message})
	}

	resp, appErr := h.service.GetByID(c.Request().Context(), rawToken, c.Param("id"))
	if appErr != nil {
		return c.JSON(appErr.Status, auth.ErrorResponse{Code: appErr.Code, Message: appErr.Message})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *TestAddressHandler) Update(c echo.Context) error {
	rawToken, ok := extractBearerToken(c.Request().Header.Get(echo.HeaderAuthorization))
	if !ok {
		appErr := auth.NewInvalidAccessToken()
		return c.JSON(appErr.Status, auth.ErrorResponse{Code: appErr.Code, Message: appErr.Message})
	}

	var req testaddress.UpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, auth.ErrorResponse{
			Code:    auth.CodeInvalidRequestBody,
			Message: "입력값 형식이 올바르지 않습니다.",
		})
	}

	resp, appErr := h.service.Update(c.Request().Context(), rawToken, c.Param("id"), req)
	if appErr != nil {
		return c.JSON(appErr.Status, auth.ErrorResponse{Code: appErr.Code, Message: appErr.Message})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *TestAddressHandler) Delete(c echo.Context) error {
	rawToken, ok := extractBearerToken(c.Request().Header.Get(echo.HeaderAuthorization))
	if !ok {
		appErr := auth.NewInvalidAccessToken()
		return c.JSON(appErr.Status, auth.ErrorResponse{Code: appErr.Code, Message: appErr.Message})
	}

	appErr := h.service.Delete(c.Request().Context(), rawToken, c.Param("id"))
	if appErr != nil {
		return c.JSON(appErr.Status, auth.ErrorResponse{Code: appErr.Code, Message: appErr.Message})
	}

	return c.NoContent(http.StatusNoContent)
}
