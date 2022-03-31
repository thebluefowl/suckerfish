package domain

import "time"

type Changelog struct {
	ID          string
	Title       string
	Description string
	State       string

	PublishedBy string
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
