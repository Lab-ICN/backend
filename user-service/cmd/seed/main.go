package main

import (
	"context"
	"flag"
	"log"

	"github.com/Lab-ICN/backend/user-service/internal/config"
	"github.com/Lab-ICN/backend/user-service/internal/postgresql"
	"github.com/Lab-ICN/backend/user-service/internal/seed"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("secret")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v\n", err)
	}
}

func main() {
	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Failed to parse config file: %v\n", err)
	}
	ctx := context.Background()
	postgresql, err := postgresql.NewConn(ctx, &cfg)
	if err != nil {
		log.Fatalf("Failed to start postgresql connection pool: %v\n", err)
	}
	flag.Parse()
	seed.Execute(ctx, postgresql, flag.Args()...)
}
