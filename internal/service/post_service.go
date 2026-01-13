package service

import (
	"context"
	"discussion-forum/internal/models"
	"discussion-forum/internal/repository"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type PostService struct {
	repo  PostRepositoryInterface
	redis *redis.Client
}

func NewPostService(repo *repository.PostRepository, redisClient *redis.Client) *PostService {
	return &PostService{
		repo:  repo,
		redis: redisClient,
	}
}

// NewPostServiceWithInterface creates a PostService with an interface (useful for testing)
func NewPostServiceWithInterface(repo PostRepositoryInterface, redisClient *redis.Client) *PostService {
	return &PostService{
		repo:  repo,
		redis: redisClient,
	}
}

func (s *PostService) Create(req *models.CreatePostRequest) (*models.Post, error) {
	post, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	s.invalidateCacheForTopic(req.TopicID)
	return post, nil
}

func (s *PostService) GetByID(id int64) (*models.Post, error) {
	return s.repo.GetByID(id)
}

func (s *PostService) GetByTopicID(topicID int64) ([]*models.Post, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("posts:topic:%d", topicID)

	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var posts []*models.Post
		if json.Unmarshal([]byte(cached), &posts) == nil {
			return posts, nil
		}
	}

	posts, err := s.repo.GetByTopicID(topicID)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(posts); err == nil {
		s.redis.Set(ctx, cacheKey, data, 3*time.Minute)
	}

	return posts, nil
}

func (s *PostService) Update(id int64, req *models.UpdatePostRequest) (*models.Post, error) {
	post, err := s.repo.Update(id, req)
	if err != nil {
		return nil, err
	}

	s.invalidateCacheForTopic(post.TopicID)
	return post, nil
}

func (s *PostService) Delete(id int64) error {
	post, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	s.invalidateCacheForTopic(post.TopicID)
	return nil
}

func (s *PostService) invalidateCacheForTopic(topicID int64) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("posts:topic:%d", topicID)
	s.redis.Del(ctx, cacheKey)
}
