package router

import (
	"github.com/labstack/echo/v4"

	"personal-mail-server/internal/http/handlers"
)

func Register(e *echo.Echo, authHandler *handlers.AuthHandler, testAddressHandler *handlers.TestAddressHandler) {
	e.GET("/healthz", func(c echo.Context) error {
		return c.NoContent(204)
	})

	e.GET("/docs", handlers.DocsPage)
	e.File("/docs/openapi.yaml", "openapi/openapi.yaml")

	v1 := e.Group("/api/v1")
	authGroup := v1.Group("/auth")
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/logout", authHandler.Logout)
	authGroup.POST("/token/reissue", authHandler.Reissue)
	testAddressGroup := v1.Group("/test-addresses")
	testAddressGroup.POST("/generate", testAddressHandler.GenerateCandidate)
}
