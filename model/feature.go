package model

import "time"

type Feature struct {
	Title       string
	Description string
	Status      string
	Owner       string
	Category    string
	Board       string
	CreatedOn   time.Time
	UpdatedAt   time.Time
}

type FeatureVote struct {
	CreatedBy string
	CreatedOn string
}

type FeatureComment struct {
	Comment   string
	CreatedBy string
	CreatedOn time.Time
	UpdatedOn time.Time
}

type FeatureCommentVote struct {
	CreatedBy string
	CreatedPm time.Time
}
