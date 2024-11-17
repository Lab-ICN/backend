package types

import "time"

type CreateUserParams struct {
	Email               string
	Username            string
	Fullname            string
	IsMember            bool
	InternshipStartDate time.Time
}
