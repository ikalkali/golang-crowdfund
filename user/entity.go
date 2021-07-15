package user

import "time"

type User struct {
	Id             int
	Name           string
	Occupation     string
	Email          string
	Role          string
	PasswordHash   string
	AvatarFileName string
	CreatedAt      time.Time
	UpdatedAt time.Time
}