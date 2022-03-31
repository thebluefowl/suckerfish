package service

import "github.com/thebluefowl/suckerfish/domain"

type UserService interface {
	CreateUser(*CreateUserRequest) *domain.User
}

type CreateUserRequest struct {
}
