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

// setupTopicTestRedis creates an in-memory Redis server for testing
func setupTopicTestRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client, mr
}

// TestTopicService_Create tests the Create method
func TestTopicService_Create(t *testing.T) {
	tests := []struct {
		name        string
		request     *models.CreateTopicRequest
		mockReturn  *models.Topic
		mockError   error
		expectError bool
	}{
		{
			name: "successful topic creation",
			request: &models.CreateTopicRequest{
				Title:  "Test Topic Title",
				UserID: 1,
			},
			mockReturn: &models.Topic{
				ID:        1,
				Title:     "Test Topic Title",
				UserID:    1,
				Username:  "testuser",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name: "repository error - invalid user",
			request: &models.CreateTopicRequest{
				Title:  "Test Topic",
				UserID: 999,
			},
			mockReturn:  nil,
			mockError:   errors.New("foreign key constraint failed"),
			expectError: true,
		},
		{
			name: "database connection error",
			request: &models.CreateTopicRequest{
				Title:  "Test Topic",
				UserID: 1,
			},
			mockReturn:  nil,
			mockError:   errors.New("database connection error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockTopicRepositoryInterface(t)
			redisClient, mr := setupTopicTestRedis(t)
			defer mr.Close()

			service := NewTopicServiceWithInterface(mockRepo, redisClient)

			mockRepo.On("Create", tt.request).Return(tt.mockReturn, tt.mockError)

			result, err := service.Create(tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockReturn.ID, result.ID)
				assert.Equal(t, tt.mockReturn.Title, result.Title)
				assert.Equal(t, tt.mockReturn.UserID, result.UserID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestTopicService_GetByID tests the GetByID method
func TestTopicService_GetByID(t *testing.T) {
	tests := []struct {
		name        string
		topicID     int64
		mockReturn  *models.Topic
		mockError   error
		expectError bool
	}{
		{
			name:    "successful retrieval",
			topicID: 1,
			mockReturn: &models.Topic{
				ID:        1,
				Title:     "Test Topic",
				UserID:    1,
				Username:  "testuser",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "topic not found",
			topicID:     999,
			mockReturn:  nil,
			mockError:   sql.ErrNoRows,
			expectError: true,
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
			mockRepo := mocks.NewMockTopicRepositoryInterface(t)
			redisClient, mr := setupTopicTestRedis(t)
			defer mr.Close()

			service := NewTopicServiceWithInterface(mockRepo, redisClient)

			mockRepo.On("GetByID", tt.topicID).Return(tt.mockReturn, tt.mockError)

			result, err := service.GetByID(tt.topicID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockReturn.ID, result.ID)
				assert.Equal(t, tt.mockReturn.Title, result.Title)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestTopicService_GetByID_WithCache tests cache hit scenario
func TestTopicService_GetByID_WithCache(t *testing.T) {
	mockRepo := mocks.NewMockTopicRepositoryInterface(t)
	redisClient, mr := setupTopicTestRedis(t)
	defer mr.Close()

	service := NewTopicServiceWithInterface(mockRepo, redisClient)

	topic := &models.Topic{
		ID:        1,
		Title:     "Cached Topic",
		UserID:    1,
		Username:  "testuser",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// First call - cache miss
	mockRepo.On("GetByID", int64(1)).Return(topic, nil).Once()

	result1, err := service.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, topic.ID, result1.ID)

	// Second call - cache hit (repository should NOT be called)
	result2, err := service.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, topic.ID, result2.ID)

	// Verify repository was only called once
	mockRepo.AssertExpectations(t)
}

// TestTopicService_GetAll tests the GetAll method
func TestTopicService_GetAll(t *testing.T) {
	tests := []struct {
		name        string
		mockReturn  []*models.Topic
		mockError   error
		expectError bool
	}{
		{
			name: "successful retrieval of topics",
			mockReturn: []*models.Topic{
				{
					ID:        1,
					Title:     "Topic 1",
					UserID:    1,
					Username:  "user1",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        2,
					Title:     "Topic 2",
					UserID:    2,
					Username:  "user2",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "empty topic list",
			mockReturn:  []*models.Topic{},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "database error",
			mockReturn:  nil,
			mockError:   errors.New("connection error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockTopicRepositoryInterface(t)
			redisClient, mr := setupTopicTestRedis(t)
			defer mr.Close()

			service := NewTopicServiceWithInterface(mockRepo, redisClient)

			mockRepo.On("GetAll").Return(tt.mockReturn, tt.mockError)

			result, err := service.GetAll()

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

// TestTopicService_GetAll_WithCache tests cache behavior
func TestTopicService_GetAll_WithCache(t *testing.T) {
	mockRepo := mocks.NewMockTopicRepositoryInterface(t)
	redisClient, mr := setupTopicTestRedis(t)
	defer mr.Close()

	service := NewTopicServiceWithInterface(mockRepo, redisClient)

	topics := []*models.Topic{
		{ID: 1, Title: "Topic 1", UserID: 1, Username: "user1"},
		{ID: 2, Title: "Topic 2", UserID: 2, Username: "user2"},
	}

	// First call - cache miss
	mockRepo.On("GetAll").Return(topics, nil).Once()

	result1, err := service.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result1))

	// Second call - cache hit (repository should NOT be called)
	result2, err := service.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result2))

	// Verify repository was only called once
	mockRepo.AssertExpectations(t)
}

// TestTopicService_Update tests the Update method
func TestTopicService_Update(t *testing.T) {
	tests := []struct {
		name        string
		topicID     int64
		request     *models.UpdateTopicRequest
		mockReturn  *models.Topic
		mockError   error
		expectError bool
	}{
		{
			name:    "successful update",
			topicID: 1,
			request: &models.UpdateTopicRequest{
				Title: "Updated Topic Title",
			},
			mockReturn: &models.Topic{
				ID:        1,
				Title:     "Updated Topic Title",
				UserID:    1,
				Username:  "testuser",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:    "topic not found",
			topicID: 999,
			request: &models.UpdateTopicRequest{
				Title: "Updated Title",
			},
			mockReturn:  nil,
			mockError:   sql.ErrNoRows,
			expectError: true,
		},
		{
			name:    "database error",
			topicID: 1,
			request: &models.UpdateTopicRequest{
				Title: "Updated Title",
			},
			mockReturn:  nil,
			mockError:   errors.New("database error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockTopicRepositoryInterface(t)
			redisClient, mr := setupTopicTestRedis(t)
			defer mr.Close()

			service := NewTopicServiceWithInterface(mockRepo, redisClient)

			mockRepo.On("Update", tt.topicID, tt.request).Return(tt.mockReturn, tt.mockError)

			result, err := service.Update(tt.topicID, tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockReturn.ID, result.ID)
				assert.Equal(t, tt.mockReturn.Title, result.Title)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestTopicService_Delete tests the Delete method
func TestTopicService_Delete(t *testing.T) {
	tests := []struct {
		name        string
		topicID     int64
		mockError   error
		expectError bool
	}{
		{
			name:        "successful deletion",
			topicID:     1,
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "topic not found",
			topicID:     999,
			mockError:   sql.ErrNoRows,
			expectError: true,
		},
		{
			name:        "foreign key constraint error",
			topicID:     1,
			mockError:   errors.New("cannot delete topic with existing posts"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockTopicRepositoryInterface(t)
			redisClient, mr := setupTopicTestRedis(t)
			defer mr.Close()

			service := NewTopicServiceWithInterface(mockRepo, redisClient)

			mockRepo.On("Delete", tt.topicID).Return(tt.mockError)

			err := service.Delete(tt.topicID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestTopicService_CacheInvalidation tests cache invalidation
func TestTopicService_CacheInvalidation(t *testing.T) {
	mockRepo := mocks.NewMockTopicRepositoryInterface(t)
	redisClient, mr := setupTopicTestRedis(t)
	defer mr.Close()

	service := NewTopicServiceWithInterface(mockRepo, redisClient)

	topics := []*models.Topic{
		{ID: 1, Title: "Topic 1", UserID: 1, Username: "user1"},
	}

	// Populate cache with GetAll
	mockRepo.On("GetAll").Return(topics, nil).Once()
	_, err := service.GetAll()
	assert.NoError(t, err)

	// Create a new topic - should invalidate cache
	newTopic := &models.Topic{
		ID:       2,
		Title:    "New Topic",
		UserID:   1,
		Username: "user1",
	}
	createReq := &models.CreateTopicRequest{
		Title:  "New Topic",
		UserID: 1,
	}

	mockRepo.On("Create", createReq).Return(newTopic, nil)
	_, err = service.Create(createReq)
	assert.NoError(t, err)

	// Next GetAll should hit repository (cache was invalidated)
	mockRepo.On("GetAll").Return(topics, nil).Once()
	_, err = service.GetAll()
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
