package app

import (
	"auth/internal/config"
	authcontroller "auth/internal/controllers/auth"
	"auth/internal/db/postgres"
	authservice "auth/internal/services/auth"
	"auth/internal/services/email"
	"auth/internal/storages"
	logutils "auth/internal/utils/log"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
)

func Run(cfg *config.Config) error {
	db, err := postgres.NewDatabase(cfg.DB)
	if err != nil {
		return errors.Wrap(err, "failed to init storage")
	}

	refreshTokenStorage := storages.NewRefreshTokenStorage(db)
	userStorage := storages.NewFakeUserStorage()

	emailService := email.NewEmailService(cfg.SMTP, userStorage)
	authService := authservice.NewAuthService(refreshTokenStorage, emailService, []byte(cfg.Auth.JWTPrivateKey), cfg.Auth.AccessTokenDuration, cfg.Auth.RefreshTokenDuration, cfg.Emails)

	authController := authcontroller.NewAuthController(authService)

	router := newRouter()
	authController.RegisterRoutes(router)

	slog.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			slog.Error("failed to start server")
		}
	}()

	slog.Info("server started")

	<-done
	slog.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "failed to stop server")
	}

	slog.Info("server stopped")

	return nil
}

func newRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logutils.NewHTTPLogger())
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	return router
}
