package domain

import "time"

type Feature struct {
	ID          string
	Title       string
	Description string
	Status      string
	Owner       string
	Category    string
	Board       string
	CreatedBy   string
	CreatedOn   time.Time
	UpdatedAt   time.Time
}

type FeatureVote struct {
	FeatureID string
	CreatedBy string
	CreatedOn string
}

type FeatureComment struct {
	FeatureID string
	Comment   string
	CreatedBy string
	CreatedOn time.Time
	UpdatedOn time.Time
}

type FeatureCommentVote struct {
	FeatureCommentID string
	CreatedBy        string
	CreatedOn        time.Time
	UpdatedOn        time.Time
}

type FeatureCommentReply struct {
	FeatureCommentID string
	Reply            string
	CreatedBy        string
	CreatedOn        time.Time
}
