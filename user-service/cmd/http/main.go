package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Lab-ICN/backend/user-service/http"
	"github.com/Lab-ICN/backend/user-service/internal/config"
	_fiber "github.com/Lab-ICN/backend/user-service/internal/fiber"
	"github.com/Lab-ICN/backend/user-service/internal/logging"
	"github.com/Lab-ICN/backend/user-service/internal/postgresql"
	"github.com/Lab-ICN/backend/user-service/repository"
	"github.com/Lab-ICN/backend/user-service/usecase"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func main() {
	content, err := os.ReadFile(os.Getenv("CONFIG_FILE"))
	if err != nil {
		log.Fatalf("Failed to open config file: %v\n", err)
	}
	cfg := new(config.Config)
	if err := json.Unmarshal(content, cfg); err != nil {
		log.Fatalf("Failed to parse config file: %v\n", err)
	}
	ctx := context.Background()

	logger := new(zap.Logger)
	if cfg.Development {
		logger, err = logging.NewDevelopment()
	} else {
		logger, err = logging.NewProduction()
	}
	if err != nil {
		log.Fatalf("Failed to build logging instance: %v\n", err)
	}
	validate := validator.New()
	postgresql, err := postgresql.NewPool(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to start postgresql connection pool: %v\n", err)
	}
	r := _fiber.New(cfg, logger)
	api := r.Group("/api")

	store := repository.NewUserPostgreSQL(postgresql)
	usecase := usecase.NewUserUsecase(store, logger)
	http.RegisterHandler(usecase, api, validate, logger)

	go func() {
		if err := r.Listen(fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)); err != nil {
			log.Panicf("Server panicked: %v\n", err)
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
		func(ctx context.Context) error {
			return logger.Sync()
		},
	)
	<-shutdownCtx.Done()
	if errors.Is(context.DeadlineExceeded, shutdownCtx.Err()) {
		log.Panicf("Shutdown timeout, force exit...")
	}
	log.Println("Gracefully shutdown...")
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
					log.Printf("Task failed: %v\n", err)
				}
			}()
		}
		wg.Wait()
		cancel()
	}()
	return _ctx
}
