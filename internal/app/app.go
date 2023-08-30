package app

import (
	"context"
	"errors"
	"github.com/zenorachi/balance-management/internal/database"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zenorachi/balance-management/pkg/auth"

	"github.com/zenorachi/balance-management/internal/repository"
	"github.com/zenorachi/balance-management/internal/service"
	"github.com/zenorachi/balance-management/internal/transport"
	"github.com/zenorachi/balance-management/pkg/hash"

	_ "github.com/lib/pq"
	"github.com/zenorachi/balance-management/internal/config"
	"github.com/zenorachi/balance-management/internal/server"
	"github.com/zenorachi/balance-management/pkg/database/postgres"
	"github.com/zenorachi/balance-management/pkg/logger"
)

const (
	shutdownTimeout = 5 * time.Second
)

func Run(cfg *config.Config) {
	err := database.DoMigrations(&cfg.DB)
	if err != nil {
		logger.Fatal("migrations", "migrations failed")
	}
	logger.Info("migrations", "migrations done")

	db, err := postgres.NewDB(&cfg.DB)
	defer func() { _ = db.Close() }()
	if err != nil {
		logger.Fatal("database-connection", err)
	}
	logger.Info("database", "postgres started")

	tokenManager := auth.NewManager(cfg.Auth.Secret)

	services := service.New(service.Deps{
		Repos:           repository.New(db),
		Hasher:          hash.NewSHA1Hasher(cfg.Auth.Salt),
		TokenManager:    tokenManager,
		AccessTokenTTL:  cfg.Auth.AccessTokenTTL,
		RefreshTokenTTL: cfg.Auth.RefreshTokenTTL,
	})

	handler := transport.NewHandler(services, tokenManager)
	srv := server.New(cfg, handler.InitRoutes())
	go func() {
		if err = srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("server", err)
		}
	}()

	logger.Info("server", "server started")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logger.Info("server", "shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		logger.Fatal("server", err)
	}
	logger.Info("server", "server stopped")
}
