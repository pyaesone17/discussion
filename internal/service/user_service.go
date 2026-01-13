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

type UserService struct {
	repo  UserRepositoryInterface
	redis *redis.Client
}

func NewUserService(repo *repository.UserRepository, redisClient *redis.Client) *UserService {
	return &UserService{
		repo:  repo,
		redis: redisClient,
	}
}

// NewUserServiceWithInterface creates a UserService with an interface (useful for testing)
func NewUserServiceWithInterface(repo UserRepositoryInterface, redisClient *redis.Client) *UserService {
	return &UserService{
		repo:  repo,
		redis: redisClient,
	}
}

func (s *UserService) Create(req *models.CreateUserRequest) (*models.User, error) {
	user, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	s.invalidateCache()
	return user, nil
}

func (s *UserService) GetByID(id int64) (*models.User, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%d", id)

	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var user models.User
		if json.Unmarshal([]byte(cached), &user) == nil {
			return &user, nil
		}
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(user); err == nil {
		s.redis.Set(ctx, cacheKey, data, 5*time.Minute)
	}

	return user, nil
}

func (s *UserService) GetAll() ([]*models.User, error) {
	ctx := context.Background()
	cacheKey := "users:all"

	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var users []*models.User
		if json.Unmarshal([]byte(cached), &users) == nil {
			return users, nil
		}
	}

	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(users); err == nil {
		s.redis.Set(ctx, cacheKey, data, 2*time.Minute)
	}

	return users, nil
}

func (s *UserService) Update(id int64, req *models.UpdateUserRequest) (*models.User, error) {
	user, err := s.repo.Update(id, req)
	if err != nil {
		return nil, err
	}

	s.invalidateCache()
	return user, nil
}

func (s *UserService) Delete(id int64) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	s.invalidateCache()
	return nil
}

func (s *UserService) invalidateCache() {
	ctx := context.Background()
	s.redis.Del(ctx, "users:all")
}
