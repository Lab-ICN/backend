package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Lab-ICN/backend/user-service/http"
	"github.com/Lab-ICN/backend/user-service/internal/config"
	_fiber "github.com/Lab-ICN/backend/user-service/internal/fiber"
	"github.com/Lab-ICN/backend/user-service/internal/postgresql"
	"github.com/Lab-ICN/backend/user-service/repository"
	"github.com/Lab-ICN/backend/user-service/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog"
)

func main() {
	content, err := os.ReadFile(os.Getenv("CONFIG_FILE"))
	if err != nil {
		stdlog.Fatalf("Failed to open config file: %v\n", err)
	}
	cfg := new(config.Config)
	if err := json.Unmarshal(content, cfg); err != nil {
		stdlog.Fatalf("Failed to parse config file: %v\n", err)
	}
	ctx := context.Background()

	logfile, err := os.Create(cfg.LogPath)
	if err != nil {
		stdlog.Fatalf("creating log file: %w", err)
	}
	log := zerolog.New(logfile).With().Timestamp().Logger()
	if !cfg.Development {
		log = log.Level(zerolog.InfoLevel)
	} else {
		log = log.Level(zerolog.DebugLevel)
	}
	validate := validator.New()
	postgresql, err := postgresql.NewPool(ctx, cfg)
	if err != nil {
		stdlog.Fatalf("Failed to start postgresql connection pool: %v\n", err)
	}
	r := _fiber.New(cfg, &log)
	r.Use(cors.New())
	api := r.Group("/backend")

	store := repository.NewUserPostgreSQL(postgresql)
	usecase := usecase.NewUserUsecase(store, &log)
	http.RegisterHandlers(usecase, cfg, api, validate)

	go func() {
		if err := r.Listen(fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)); err != nil {
			stdlog.Panicf("Server panicked: %v\n", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-sig

	shutdownCtx := gracefulShutdown(ctx, 5*time.Second,
		func(ctx context.Context) error {
			return r.Shutdown()
		},
		func(ctx context.Context) error {
			postgresql.Close()
			return nil
		},
	)
	<-shutdownCtx.Done()
	if errors.Is(context.DeadlineExceeded, shutdownCtx.Err()) {
		stdlog.Panicf("Shutdown timeout, force exit...")
	}
	stdlog.Println("Gracefully shutdown...")
}

func gracefulShutdown(
	ctx context.Context,
	timeout time.Duration,
	tasks ...func(ctx context.Context) error,
) context.Context {
	_ctx, cancel := context.WithTimeout(ctx, timeout)
	go func() {
		var wg sync.WaitGroup
		for _, task := range tasks {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := task(ctx); err != nil {
					stdlog.Printf("Task failed: %v\n", err)
				}
			}()
		}
		wg.Wait()
		cancel()
	}()
	return _ctx
}
