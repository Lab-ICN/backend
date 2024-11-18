package types

import "time"

type User struct {
	ID                  uint64    `json:"id"`
	Email               string    `json:"email"`
	Username            string    `json:"username"`
	Fullname            string    `json:"fullname"`
	IsMember            bool      `json:"isMember"`
	InternshipStartDate time.Time `json:"internshipStartDate"`
}

type CreateUserParams struct {
	Email               string
	Username            string
	Fullname            string
	IsMember            bool
	InternshipStartDate time.Time
}
