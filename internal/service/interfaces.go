package service

import "discussion-forum/internal/models"

// Repository interfaces for dependency injection and testing

type UserRepositoryInterface interface {
	Create(req *models.CreateUserRequest) (*models.User, error)
	GetByID(id int64) (*models.User, error)
	GetAll() ([]*models.User, error)
	Update(id int64, req *models.UpdateUserRequest) (*models.User, error)
	Delete(id int64) error
}

type TopicRepositoryInterface interface {
	Create(req *models.CreateTopicRequest) (*models.Topic, error)
	GetByID(id int64) (*models.Topic, error)
	GetAll() ([]*models.Topic, error)
	Update(id int64, req *models.UpdateTopicRequest) (*models.Topic, error)
	Delete(id int64) error
	Search(query string) ([]*models.Topic, error)
}

type PostRepositoryInterface interface {
	Create(req *models.CreatePostRequest) (*models.Post, error)
	GetByID(id int64) (*models.Post, error)
	GetByTopicID(topicID int64) ([]*models.Post, error)
	Update(id int64, req *models.UpdatePostRequest) (*models.Post, error)
	Delete(id int64) error
}
