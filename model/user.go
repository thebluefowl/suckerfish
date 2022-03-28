package model

type User struct {
	Name   string
	Avatar string
	Email  string
}

type Account struct {
	Email    string
	Provider string
	Token    string
}
