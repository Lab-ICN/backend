package repository

import (
	"time"
)

type User struct {
	ID                  int64
	Email               string
	Username            string
	Fullname            string
	IsMember            bool
	InternshipStartDate time.Time
}
