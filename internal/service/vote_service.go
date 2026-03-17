package service

import (
	"context"
	"discussion-forum/internal/models"
	"discussion-forum/internal/repository"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type VoteService struct {
	voteRepo VoteRepositoryInterface
	postRepo PostRepositoryInterface
	redis    *redis.Client
}

func NewVoteService(voteRepo *repository.VoteRepository, postRepo *repository.PostRepository, redisClient *redis.Client) *VoteService {
	return &VoteService{
		voteRepo: voteRepo,
		postRepo: postRepo,
		redis:    redisClient,
	}
}

func (s *VoteService) CastVote(userID, votableID int64, votableType string, value int) error {
	if err := s.voteRepo.Upsert(userID, votableID, votableType, value); err != nil {
		return err
	}

	s.invalidateCache(votableID, votableType)
	return nil
}

func (s *VoteService) RemoveVote(userID, votableID int64, votableType string) error {
	if err := s.voteRepo.Delete(userID, votableID, votableType); err != nil {
		return err
	}

	s.invalidateCache(votableID, votableType)
	return nil
}

func (s *VoteService) invalidateCache(votableID int64, votableType string) {
	ctx := context.Background()

	switch votableType {
	case models.VotableTypeTopic:
		s.redis.Del(ctx, fmt.Sprintf("topic:%d", votableID))
		s.redis.Del(ctx, "topics:all")
	case models.VotableTypePost:
		s.redis.Del(ctx, fmt.Sprintf("post:%d", votableID))
		// Also invalidate the post's topic list cache
		post, err := s.postRepo.GetByID(votableID)
		if err == nil {
			s.redis.Del(ctx, fmt.Sprintf("posts:topic:%d", post.TopicID))
		}
	}
}
