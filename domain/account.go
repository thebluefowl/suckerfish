package domain

type Account struct {
	Email    string
	Provider string
	Token    string

	ID     string
	Name   string
	Team   string
	Avatar string
}
