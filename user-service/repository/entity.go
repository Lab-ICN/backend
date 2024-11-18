package repository

import (
	"time"

	"github.com/Lab-ICN/backend/user-service/types"
)

type User struct {
	ID                  uint64
	Email               string
	Username            string
	Fullname            string
	IsMember            bool
	InternshipStartDate time.Time
}

func (u User) DTO() types.User {
	return types.User{
		ID:                  u.ID,
		Email:               u.Email,
		Username:            u.Username,
		Fullname:            u.Fullname,
		IsMember:            u.IsMember,
		InternshipStartDate: u.InternshipStartDate,
	}
}
