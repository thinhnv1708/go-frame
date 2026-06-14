package entity

import "time"

type User struct {
	ID       string
	Name     string
	Username string
	Password string
	Dob      time.Time
}
