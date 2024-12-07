package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/Lab-ICN/backend/user-service/internal/config"
	"github.com/Lab-ICN/backend/user-service/internal/postgresql"
	"github.com/Lab-ICN/backend/user-service/internal/seed"
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
	postgresql, err := postgresql.NewConn(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to start postgresql connection pool: %v\n", err)
	}
	flag.Parse()
	seed.Execute(ctx, postgresql, flag.Args()...)
}
