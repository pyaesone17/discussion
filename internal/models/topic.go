package models

import "time"

type Topic struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Score     int64     `json:"score"`
}

type CreateTopicRequest struct {
	Title  string `json:"title" binding:"required,min=5,max=200"`
	UserID int64  `json:"user_id" binding:"required"`
}

type UpdateTopicRequest struct {
	Title string `json:"title" binding:"omitempty,min=5,max=200"`
}
