package domain

import "time"

type User struct {
	ID        string
	Name      string
	AvatarURL string
	Email     string
	Provider  string
	Token     string
	Company   string
	Location  string
	IsStaff   bool
	IsNewUser bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepository interface {
	GetByEmail(email string) (*User, error)
	Create(*User) error
}
