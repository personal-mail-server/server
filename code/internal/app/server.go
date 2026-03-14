package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"personal-mail-server/internal/auth"
	"personal-mail-server/internal/config"
	"personal-mail-server/internal/db"
	"personal-mail-server/internal/http/handlers"
	"personal-mail-server/internal/http/router"
)

type Server struct {
	echo       *echo.Echo
	httpServer *http.Server
	cleanup    func()
}

func NewServer(cfg config.Config) (*Server, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.RunMigrations(ctx, pool, "migrations"); err != nil {
		pool.Close()
		return nil, err
	}

	issuer := auth.NewJWTIssuer(cfg.AccessTokenSecret, cfg.RefreshTokenSecret)
	repo := auth.NewPostgresRepository(pool)
	service := auth.NewService(repo, issuer, auth.RealClock{})
	authHandler := handlers.NewAuthHandler(service)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("1M"))

	if cfg.EnableCORS {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: cfg.AllowedOrigins,
			AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		}))
	}

	router.Register(e, authHandler)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           e,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &Server{echo: e, httpServer: server, cleanup: pool.Close}, nil
}

func (s *Server) Start() error {
	return s.echo.StartServer(s.httpServer)
}

func (s *Server) Shutdown(ctx context.Context) error {
	defer func() {
		if s.cleanup != nil {
			s.cleanup()
		}
	}()
	return s.echo.Shutdown(ctx)
}
