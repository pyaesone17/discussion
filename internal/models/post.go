package models

import "time"

type Post struct {
	ID        int64     `json:"id"`
	TopicID   int64     `json:"topic_id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePostRequest struct {
	TopicID int64  `json:"topic_id" binding:"required"`
	UserID  int64  `json:"user_id" binding:"required"`
	Content string `json:"content" binding:"required,min=1,max=5000"`
}

type UpdatePostRequest struct {
	Content string `json:"content" binding:"omitempty,min=1,max=5000"`
}
