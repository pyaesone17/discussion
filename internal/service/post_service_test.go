package service

import (
	"database/sql"
	"discussion-forum/internal/models"
	"discussion-forum/internal/service/mocks"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// setupPostTestRedis creates an in-memory Redis server for testing
func setupPostTestRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client, mr
}

// TestPostService_Create tests the Create method
func TestPostService_Create(t *testing.T) {
	tests := []struct {
		name        string
		request     *models.CreatePostRequest
		mockReturn  *models.Post
		mockError   error
		expectError bool
	}{
		{
			name: "successful post creation",
			request: &models.CreatePostRequest{
				TopicID: 1,
				UserID:  1,
				Content: "This is a test post content",
			},
			mockReturn: &models.Post{
				ID:        1,
				TopicID:   1,
				UserID:    1,
				Username:  "testuser",
				Content:   "This is a test post content",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name: "invalid topic ID",
			request: &models.CreatePostRequest{
				TopicID: 999,
				UserID:  1,
				Content: "Test content",
			},
			mockReturn:  nil,
			mockError:   errors.New("foreign key constraint failed - topic does not exist"),
			expectError: true,
		},
		{
			name: "invalid user ID",
			request: &models.CreatePostRequest{
				TopicID: 1,
				UserID:  999,
				Content: "Test content",
			},
			mockReturn:  nil,
			mockError:   errors.New("foreign key constraint failed - user does not exist"),
			expectError: true,
		},
		{
			name: "database connection error",
			request: &models.CreatePostRequest{
				TopicID: 1,
				UserID:  1,
				Content: "Test content",
			},
			mockReturn:  nil,
			mockError:   errors.New("database connection error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockPostRepositoryInterface(t)
			redisClient, mr := setupPostTestRedis(t)
			defer mr.Close()

			service := NewPostServiceWithInterface(mockRepo, redisClient)

			mockRepo.On("Create", tt.request).Return(tt.mockReturn, tt.mockError)

			result, err := service.Create(tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockReturn.ID, result.ID)
				assert.Equal(t, tt.mockReturn.Content, result.Content)
				assert.Equal(t, tt.mockReturn.TopicID, result.TopicID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestPostService_GetByID tests the GetByID method
func TestPostService_GetByID(t *testing.T) {
	tests := []struct {
		name        string
		postID      int64
		mockReturn  *models.Post
		mockError   error
		expectError bool
	}{
		{
			name:   "successful retrieval",
			postID: 1,
			mockReturn: &models.Post{
				ID:        1,
				TopicID:   1,
				UserID:    1,
				Username:  "testuser",
				Content:   "Test post content",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "post not found",
			postID:      999,
			mockReturn:  nil,
			mockError:   sql.ErrNoRows,
			expectError: true,
		},
		{
			name:        "database error",
			postID:      1,
			mockReturn:  nil,
			mockError:   errors.New("connection error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockPostRepositoryInterface(t)
			redisClient, mr := setupPostTestRedis(t)
			defer mr.Close()

			service := NewPostServiceWithInterface(mockRepo, redisClient)

			mockRepo.On("GetByID", tt.postID).Return(tt.mockReturn, tt.mockError)

			result, err := service.GetByID(tt.postID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockReturn.ID, result.ID)
				assert.Equal(t, tt.mockReturn.Content, result.Content)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestPostService_GetByTopicID tests the GetByTopicID method
func TestPostService_GetByTopicID(t *testing.T) {
	tests := []struct {
		name        string
		topicID     int64
		mockReturn  []*models.Post
		mockError   error
		expectError bool
	}{
		{
			name:    "successful retrieval of posts",
			topicID: 1,
			mockReturn: []*models.Post{
				{
					ID:        1,
					TopicID:   1,
					UserID:    1,
					Username:  "user1",
					Content:   "First post",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        2,
					TopicID:   1,
					UserID:    2,
					Username:  "user2",
					Content:   "Second post",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "empty post list",
			topicID:     1,
			mockReturn:  []*models.Post{},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "database error",
			topicID:     1,
			mockReturn:  nil,
			mockError:   errors.New("connection error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockPostRepositoryInterface(t)
			redisClient, mr := setupPostTestRedis(t)
			defer mr.Close()

			service := NewPostServiceWithInterface(mockRepo, redisClient)

			mockRepo.On("GetByTopicID", tt.topicID).Return(tt.mockReturn, tt.mockError)

			result, err := service.GetByTopicID(tt.topicID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, len(tt.mockReturn), len(result))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestPostService_GetByTopicID_WithCache tests cache hit scenario
func TestPostService_GetByTopicID_WithCache(t *testing.T) {
	mockRepo := mocks.NewMockPostRepositoryInterface(t)
	redisClient, mr := setupPostTestRedis(t)
	defer mr.Close()

	service := NewPostServiceWithInterface(mockRepo, redisClient)

	posts := []*models.Post{
		{ID: 1, TopicID: 1, UserID: 1, Username: "testuser", Content: "Test post"},
	}

	// First call - cache miss
	mockRepo.On("GetByTopicID", int64(1)).Return(posts, nil).Once()

	result1, err := service.GetByTopicID(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result1))

	// Second call - cache hit (repository should NOT be called)
	result2, err := service.GetByTopicID(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result2))

	// Verify repository was only called once
	mockRepo.AssertExpectations(t)
}

// TestPostService_Update tests the Update method
func TestPostService_Update(t *testing.T) {
	tests := []struct {
		name        string
		postID      int64
		request     *models.UpdatePostRequest
		mockReturn  *models.Post
		mockError   error
		expectError bool
	}{
		{
			name:   "successful update",
			postID: 1,
			request: &models.UpdatePostRequest{
				Content: "Updated post content",
			},
			mockReturn: &models.Post{
				ID:        1,
				TopicID:   1,
				UserID:    1,
				Username:  "testuser",
				Content:   "Updated post content",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:   "post not found",
			postID: 999,
			request: &models.UpdatePostRequest{
				Content: "Updated content",
			},
			mockReturn:  nil,
			mockError:   sql.ErrNoRows,
			expectError: true,
		},
		{
			name:   "database error",
			postID: 1,
			request: &models.UpdatePostRequest{
				Content: "Updated content",
			},
			mockReturn:  nil,
			mockError:   errors.New("database error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockPostRepositoryInterface(t)
			redisClient, mr := setupPostTestRedis(t)
			defer mr.Close()

			service := NewPostServiceWithInterface(mockRepo, redisClient)

			mockRepo.On("Update", tt.postID, tt.request).Return(tt.mockReturn, tt.mockError)

			result, err := service.Update(tt.postID, tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockReturn.ID, result.ID)
				assert.Equal(t, tt.mockReturn.Content, result.Content)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestPostService_Delete tests the Delete method
func TestPostService_Delete(t *testing.T) {
	tests := []struct {
		name         string
		postID       int64
		mockGetPost  *models.Post
		mockGetError error
		mockDelError error
		expectError  bool
	}{
		{
			name:   "successful deletion",
			postID: 1,
			mockGetPost: &models.Post{
				ID:       1,
				TopicID:  1,
				UserID:   1,
				Username: "testuser",
				Content:  "Test content",
			},
			mockGetError: nil,
			mockDelError: nil,
			expectError:  false,
		},
		{
			name:         "post not found on get",
			postID:       999,
			mockGetPost:  nil,
			mockGetError: sql.ErrNoRows,
			mockDelError: nil,
			expectError:  true,
		},
		{
			name:   "deletion error",
			postID: 1,
			mockGetPost: &models.Post{
				ID:       1,
				TopicID:  1,
				UserID:   1,
				Username: "testuser",
				Content:  "Test content",
			},
			mockGetError: nil,
			mockDelError: errors.New("database error"),
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockPostRepositoryInterface(t)
			redisClient, mr := setupPostTestRedis(t)
			defer mr.Close()

			service := NewPostServiceWithInterface(mockRepo, redisClient)

			mockRepo.On("GetByID", tt.postID).Return(tt.mockGetPost, tt.mockGetError)
			if tt.mockGetError == nil {
				mockRepo.On("Delete", tt.postID).Return(tt.mockDelError)
			}

			err := service.Delete(tt.postID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestPostService_CacheInvalidation tests cache invalidation on create
func TestPostService_CacheInvalidation(t *testing.T) {
	mockRepo := mocks.NewMockPostRepositoryInterface(t)
	redisClient, mr := setupPostTestRedis(t)
	defer mr.Close()

	service := NewPostServiceWithInterface(mockRepo, redisClient)

	posts := []*models.Post{
		{ID: 1, TopicID: 1, UserID: 1, Username: "user1", Content: "Post 1"},
	}

	// Populate cache with GetByTopicID
	mockRepo.On("GetByTopicID", int64(1)).Return(posts, nil).Once()
	_, err := service.GetByTopicID(1)
	assert.NoError(t, err)

	// Create a new post - should invalidate cache for this topic
	newPost := &models.Post{
		ID:       2,
		TopicID:  1,
		UserID:   1,
		Username: "user1",
		Content:  "New post",
	}
	createReq := &models.CreatePostRequest{
		TopicID: 1,
		UserID:  1,
		Content: "New post",
	}

	mockRepo.On("Create", createReq).Return(newPost, nil)
	_, err = service.Create(createReq)
	assert.NoError(t, err)

	// Next GetByTopicID should hit repository (cache was invalidated)
	mockRepo.On("GetByTopicID", int64(1)).Return(posts, nil).Once()
	_, err = service.GetByTopicID(1)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
