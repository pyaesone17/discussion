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

type TopicService struct {
	repo  TopicRepositoryInterface
	redis *redis.Client
}

func NewTopicService(repo *repository.TopicRepository, redisClient *redis.Client) *TopicService {
	return &TopicService{
		repo:  repo,
		redis: redisClient,
	}
}

// NewTopicServiceWithInterface creates a TopicService with an interface (useful for testing)
func NewTopicServiceWithInterface(repo TopicRepositoryInterface, redisClient *redis.Client) *TopicService {
	return &TopicService{
		repo:  repo,
		redis: redisClient,
	}
}

func (s *TopicService) Create(req *models.CreateTopicRequest) (*models.Topic, error) {
	topic, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	s.invalidateCache()
	return topic, nil
}

func (s *TopicService) GetByID(id int64) (*models.Topic, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("topic:%d", id)

	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var topic models.Topic
		if json.Unmarshal([]byte(cached), &topic) == nil {
			return &topic, nil
		}
	}

	topic, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(topic); err == nil {
		s.redis.Set(ctx, cacheKey, data, 5*time.Minute)
	}

	return topic, nil
}

func (s *TopicService) GetAll() ([]*models.Topic, error) {
	ctx := context.Background()
	cacheKey := "topics:all"

	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var topics []*models.Topic
		if json.Unmarshal([]byte(cached), &topics) == nil {
			return topics, nil
		}
	}

	topics, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(topics); err == nil {
		s.redis.Set(ctx, cacheKey, data, 2*time.Minute)
	}

	return topics, nil
}

func (s *TopicService) Update(id int64, req *models.UpdateTopicRequest) (*models.Topic, error) {
	topic, err := s.repo.Update(id, req)
	if err != nil {
		return nil, err
	}

	s.invalidateCache()
	return topic, nil
}

func (s *TopicService) Delete(id int64) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	s.invalidateCache()
	return nil
}

func (s *TopicService) invalidateCache() {
	ctx := context.Background()
	s.redis.Del(ctx, "topics:all")
}
