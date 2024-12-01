package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"auth-service/config"
	"auth-service/internal/app/authservice"
	"auth-service/internal/bootstrap"
	"auth-service/internal/http/handlers"
	"auth-service/internal/repository/tokenrepository"
	"auth-service/internal/repository/userrepository"
)

func Run(cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())

	pgConnPool, err := bootstrap.InitDB(ctx, cfg)
	if err != nil {
		logrus.Fatalf("failed to connect postgres %s, %v", cfg.PgDSN, err)
	}

	userRepository := userrepository.NewUserRepository(pgConnPool)
	tokenRepository := tokenrepository.NewTokenRepository(pgConnPool)

	authService := authservice.NewAuthService(
		cfg,
		userRepository,
		tokenRepository,
	)

	handler := handlers.NewHandler(authService)

	s := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        handler.InitRoutes(),
		MaxHeaderBytes: 1 << 28, // 1MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	go func() {
		logrus.Infof("Starting listening http server at %s", cfg.Port)

		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("error service http server %v", err)
		}
	}()

	gracefulShotdown(ctx, s, cancel)

	return nil
}

func gracefulShotdown(ctx context.Context, s *http.Server, cancel context.CancelFunc) {
	const waitTime = 5 * time.Second

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(ch)

	sig := <-ch
	logrus.Infof("Received shutdown signal: %v. Initiating graceful shutdown...", sig)

	if err := s.Shutdown(ctx); err != nil {
		logrus.Errorf("error shutting down server: %v", err)
	}

	cancel()
	time.Sleep(waitTime)
	logrus.Info("Graceful shutdown completed.")
}
