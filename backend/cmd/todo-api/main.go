// Command todo-api is the HTTP server of the TODO application.
// main stays thin: parse config, wire dependencies, call run, exit.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/benja/todo-api/internal/config"
	"github.com/benja/todo-api/internal/handler"
	"github.com/benja/todo-api/internal/service"
	"github.com/benja/todo-api/internal/store/postgres"
)

func main() {
	if err := run(context.Background(), os.Getenv); err != nil {
		slog.Error("server stopped", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(ctx context.Context, getenv func(string) string) error {
	// .env is a local development convenience; production provides real
	// environment variables. A missing file is not an error.
	_ = godotenv.Load()

	cfg, err := config.Load(getenv)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Ctrl+C / SIGTERM starts graceful shutdown.
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := postgres.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer pool.Close()

	todoSvc := service.NewTodoService(postgres.NewTodoStore(pool))

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler.NewRouter(cfg.FrontendURL, todoSvc, pool),
		ReadHeaderTimeout: 5 * time.Second,   // Slowloris guard
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	errc := make(chan error, 1)
	go func() {
		errc <- srv.ListenAndServe()
	}()
	slog.Info("listening", slog.String("addr", cfg.Addr))

	select {
	case err := <-errc:
		return fmt.Errorf("listen: %w", err) // e.g. port already in use
	case <-ctx.Done():
	}

	// ctx is already cancelled; shutdown needs its own deadline.
	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	if err := <-errc; !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
