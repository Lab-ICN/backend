package main

import (
	"context"
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
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("secret")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v\n", err)
	}
	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Failed to parse config file: %v\n", err)
	}
	ctx := context.Background()

	logging, err := logging.New(&cfg)
	if err != nil {
		log.Fatalf("Failed to build logging instance: %v\n", err)
	}
	validate := validator.New()
	r := _fiber.New(&cfg)
	postgresql, err := postgresql.NewPool(ctx, &cfg)
	if err != nil {
		log.Fatalf("Failed to start postgresql connection pool: %v\n", err)
	}

	store := repository.NewUserPostgreSQL(postgresql)
	usecase := usecase.NewUserUsecase(store, logging)
	http.RegisterHandler(usecase, r, validate, logging)

	go func() {
		if err := r.Listen(fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)); err != nil {
			log.Panicf("Server panicked: %v\n", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-sig

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	go gracefulShutdown(timeoutCtx, cancel,
		func(ctx context.Context) error {
			return r.Shutdown()
		},
		func(ctx context.Context) error {
			postgresql.Close()
			return nil
		},
		func(ctx context.Context) error {
			return logging.Sync()
		},
	)
	<-timeoutCtx.Done()
	if errors.Is(context.DeadlineExceeded, timeoutCtx.Err()) {
		log.Panicf("Shutdown timeout, force exit...")
	}
	log.Println("Gracefully shutdown...")
}

func gracefulShutdown(
	ctx context.Context,
	cancel context.CancelFunc,
	tasks ...func(ctx context.Context) error,
) {
	defer cancel()
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
}
