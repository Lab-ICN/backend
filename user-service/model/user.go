package model

import "time"

type User struct {
	ID                    int64
	email                 string
	username              string
	fullname              string
	is_member             bool
	internship_start_date time.Time
}
