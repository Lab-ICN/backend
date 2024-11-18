package seed

import (
	"context"
	"log"
	"math/rand"
	"reflect"
	"time"

	"github.com/Lab-ICN/backend/user-service/repository"
	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5"
)

const (
	size = 1000
)

type seed struct{ conn *pgx.Conn }

func Execute(ctx context.Context, conn *pgx.Conn, methods ...string) {
	seed := seed{conn}
	log.Println("Running seeder...")
	for _, method := range methods {
		reflect.ValueOf(&seed).MethodByName(method).Call([]reflect.Value{reflect.ValueOf(ctx)})
	}
	log.Println("Finised running seeder...")
}

func (s *seed) Users(ctx context.Context) {
	users := make([]repository.User, size)
	for i := range size {
		timestamp, _ := time.Parse(time.DateTime, faker.Timestamp())
		users[i] = repository.User{
			Email:               faker.Email(),
			Username:            faker.Username(),
			Fullname:            faker.Name(),
			IsMember:            rand.Intn(2) == 1,
			InternshipStartDate: timestamp,
		}
	}
	s.conn.CopyFrom(
		ctx,
		pgx.Identifier{"users"},
		[]string{"email", "username", "fullname", "is_member", "internship_start_date"},
		pgx.CopyFromSlice(len(users), func(i int) ([]interface{}, error) {
			return []interface{}{
				users[i].Email,
				users[i].Username,
				users[i].Fullname,
				users[i].IsMember,
				users[i].InternshipStartDate,
			}, nil
		}),
	)
}
