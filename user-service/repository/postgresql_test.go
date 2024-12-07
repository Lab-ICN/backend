package repository_test

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Lab-ICN/backend/user-service/internal/config"
	"github.com/Lab-ICN/backend/user-service/internal/postgresql"
	"github.com/Lab-ICN/backend/user-service/repository"
	"github.com/Lab-ICN/backend/user-service/types"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

var (
	conn  *pgxpool.Pool
	store repository.IUserStorage
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	var err error
	content, err := os.ReadFile(os.Getenv("CONFIG_FILE"))
	if err != nil {
		log.Fatalf("Failed to open config file: %v\n", err)
	}
	cfg := new(config.Config)
	if err := json.Unmarshal(content, cfg); err != nil {
		log.Fatalf("Failed to parse config file: %v\n", err)
	}
	conn, err = postgresql.NewPool(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to start postgresql connection pool: %v\n", err)
	}
	defer conn.Close()
	store = repository.NewUserPostgreSQL(conn)
	code := m.Run()
	_, err = conn.Exec(ctx, `DELETE FROM users`)
	if err != nil {
		log.Fatalf("Failed to do cleanup task: %v\n", err)
	}
	os.Exit(code)
}

func TestCreate(t *testing.T) {
	ctx := context.Background()
	user := &types.CreateUserParams{
		Email:               "test@example.com",
		Username:            "testuser",
		Fullname:            "Test User",
		IsMember:            true,
		InternshipStartDate: time.Now(),
	}
	id, err := store.Create(ctx, user)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, id)
}

func TestCreateBulk(t *testing.T) {
	ctx := context.Background()
	users := []types.CreateUserParams{
		{
			Email:               "jane.doe@example.com",
			Username:            "janedoe123",
			Fullname:            "Jane Doe",
			IsMember:            true,
			InternshipStartDate: time.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			Email:               "john.smith@example.com",
			Username:            "johnsmith98",
			Fullname:            "John Smith",
			IsMember:            false,
			InternshipStartDate: time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Email:               "alice.wonderland@example.com",
			Username:            "alicew",
			Fullname:            "Alice Wonderland",
			IsMember:            true,
			InternshipStartDate: time.Date(2024, time.February, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			Email:               "bob.builder@example.com",
			Username:            "builderbob",
			Fullname:            "Bob Builder",
			IsMember:            false,
			InternshipStartDate: time.Date(2024, time.April, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			Email:               "charlie.brown@example.com",
			Username:            "charlieb",
			Fullname:            "Charlie Brown",
			IsMember:            true,
			InternshipStartDate: time.Date(2024, time.May, 5, 0, 0, 0, 0, time.UTC),
		},
	}
	err := store.CreateBulk(ctx, users)
	assert.Nil(t, err)
}

func TestList(t *testing.T) {
	ctx := context.Background()
	_, err := store.List(ctx)
	assert.Nil(t, err)
}

func TestListPassed(t *testing.T) {
	ctx := context.Background()
	_, err := store.ListPassed(ctx, 2024)
	assert.Nil(t, err)
}
