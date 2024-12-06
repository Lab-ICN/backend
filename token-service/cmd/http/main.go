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

	"github.com/Lab-ICN/backend/token-service/internal/config"
	_fiber "github.com/Lab-ICN/backend/token-service/internal/fiber"
	"github.com/Lab-ICN/backend/token-service/internal/http"
	"github.com/Lab-ICN/backend/token-service/internal/postgresql"
	"github.com/Lab-ICN/backend/token-service/internal/repository"
	"github.com/Lab-ICN/backend/token-service/internal/usecase"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog"
)

func main() {
	path := os.Getenv("CONFIG_FILE")
	content, err := os.ReadFile(path)
	if err != nil {
		stdlog.Fatalf("opening config file at %s: %v\n", path, err)
	}
	cfg := new(config.Config)
	if err := json.Unmarshal(content, cfg); err != nil {
		stdlog.Fatalf("parsing config file content: %v\n", err)
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
		stdlog.Fatalf("starting postgresql connection pool: %v\n", err)
	}
	r := _fiber.New(cfg, &log)
	r.Use(cors.New())
	api := r.Group("/backend")

	repo := repository.NewTokenPostgreSQL(postgresql)
	usecase := usecase.NewTokenUsecase(repo, cfg)
	http.RegisterHandlers(usecase, cfg, api, validate)

	go func() {
		if err := r.Listen(fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)); err != nil {
			stdlog.Panicf("server listening connection: %v\n", err)
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
		stdlog.Panicf("graceful shutdown timed out, force exiting...")
	}
	stdlog.Println("gracefully shutdown")
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
					stdlog.Printf("executing shutdown task: %v\n", err)
				}
			}()
		}
		wg.Wait()
		cancel()
	}()
	return _ctx
}
