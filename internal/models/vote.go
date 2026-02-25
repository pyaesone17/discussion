package models

import "time"

const (
	VotableTypeTopic = "topic"
	VotableTypePost  = "post"
)

type Vote struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	VotableType string    `json:"votable_type"`
	VotableID   int64     `json:"votable_id"`
	Value       int       `json:"value"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type VoteRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
	Value  int   `json:"value"`
}
