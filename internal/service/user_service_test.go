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

// setupTestRedis creates an in-memory Redis server for testing
func setupTestRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client, mr
}

// TestUserService_Create tests the Create method
func TestUserService_Create(t *testing.T) {
	tests := []struct {
		name        string
		request     *models.CreateUserRequest
		mockReturn  *models.User
		mockError   error
		expectError bool
	}{
		{
			name: "successful user creation",
			request: &models.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
			},
			mockReturn: &models.User{
				ID:        1,
				Username:  "testuser",
				Email:     "test@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name: "repository error",
			request: &models.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
			},
			mockReturn:  nil,
			mockError:   errors.New("database error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockUserRepositoryInterface(t)
			redisClient, mr := setupTestRedis(t)
			defer mr.Close()

			service := NewUserServiceWithInterface(mockRepo, redisClient)

			mockRepo.On("Create", tt.request).Return(tt.mockReturn, tt.mockError)

			result, err := service.Create(tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockReturn.ID, result.ID)
				assert.Equal(t, tt.mockReturn.Username, result.Username)
				assert.Equal(t, tt.mockReturn.Email, result.Email)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestUserService_GetByID tests the GetByID method
func TestUserService_GetByID(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		mockReturn  *models.User
		mockError   error
		expectError bool
	}{
		{
			name:   "successful retrieval - cache miss",
			userID: 1,
			mockReturn: &models.User{
				ID:        1,
				Username:  "testuser",
				Email:     "test@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "user not found",
			userID:      999,
			mockReturn:  nil,
			mockError:   sql.ErrNoRows,
			expectError: true,
		},
		{
			name:        "database error",
			userID:      1,
			mockReturn:  nil,
			mockError:   errors.New("connection error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockUserRepositoryInterface(t)
			redisClient, mr := setupTestRedis(t)
			defer mr.Close()

			service := NewUserServiceWithInterface(mockRepo, redisClient)

			mockRepo.On("GetByID", tt.userID).Return(tt.mockReturn, tt.mockError)

			result, err := service.GetByID(tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.userID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestUserService_GetByID_WithCache tests cache hit scenario
func TestUserService_GetByID_WithCache(t *testing.T) {
	mockRepo := mocks.NewMockUserRepositoryInterface(t)
	redisClient, mr := setupTestRedis(t)
	defer mr.Close()

	service := NewUserServiceWithInterface(mockRepo, redisClient)

	user := &models.User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// First call - cache miss
	mockRepo.On("GetByID", int64(1)).Return(user, nil).Once()

	result1, err := service.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, result1.ID)

	// Second call - cache hit (repository should NOT be called)
	result2, err := service.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, result2.ID)

	// Verify repository was only called once
	mockRepo.AssertExpectations(t)
}

// TestUserService_GetAll tests the GetAll method
func TestUserService_GetAll(t *testing.T) {
	tests := []struct {
		name        string
		mockReturn  []*models.User
		mockError   error
		expectError bool
	}{
		{
			name: "successful retrieval of users",
			mockReturn: []*models.User{
				{
					ID:        1,
					Username:  "user1",
					Email:     "user1@example.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        2,
					Username:  "user2",
					Email:     "user2@example.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "empty user list",
			mockReturn:  []*models.User{},
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
			mockRepo := mocks.NewMockUserRepositoryInterface(t)
			redisClient, mr := setupTestRedis(t)
			defer mr.Close()

			service := NewUserServiceWithInterface(mockRepo, redisClient)

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

// TestUserService_GetAll_WithCache tests cache behavior for GetAll
func TestUserService_GetAll_WithCache(t *testing.T) {
	mockRepo := mocks.NewMockUserRepositoryInterface(t)
	redisClient, mr := setupTestRedis(t)
	defer mr.Close()

	service := NewUserServiceWithInterface(mockRepo, redisClient)

	users := []*models.User{
		{ID: 1, Username: "user1", Email: "user1@example.com"},
		{ID: 2, Username: "user2", Email: "user2@example.com"},
	}

	// First call - cache miss
	mockRepo.On("GetAll").Return(users, nil).Once()

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

// TestUserService_Update tests the Update method
func TestUserService_Update(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		request     *models.UpdateUserRequest
		mockReturn  *models.User
		mockError   error
		expectError bool
	}{
		{
			name:   "successful update",
			userID: 1,
			request: &models.UpdateUserRequest{
				Username: "updateduser",
				Email:    "updated@example.com",
			},
			mockReturn: &models.User{
				ID:        1,
				Username:  "updateduser",
				Email:     "updated@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:   "update username only",
			userID: 1,
			request: &models.UpdateUserRequest{
				Username: "newusername",
			},
			mockReturn: &models.User{
				ID:        1,
				Username:  "newusername",
				Email:     "test@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:   "user not found",
			userID: 999,
			request: &models.UpdateUserRequest{
				Username: "test",
			},
			mockReturn:  nil,
			mockError:   sql.ErrNoRows,
			expectError: true,
		},
		{
			name:   "database error",
			userID: 1,
			request: &models.UpdateUserRequest{
				Email: "test@example.com",
			},
			mockReturn:  nil,
			mockError:   errors.New("database error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockUserRepositoryInterface(t)
			redisClient, mr := setupTestRedis(t)
			defer mr.Close()

			service := NewUserServiceWithInterface(mockRepo, redisClient)

			mockRepo.On("Update", tt.userID, tt.request).Return(tt.mockReturn, tt.mockError)

			result, err := service.Update(tt.userID, tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockReturn.ID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestUserService_Delete tests the Delete method
func TestUserService_Delete(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		mockError   error
		expectError bool
	}{
		{
			name:        "successful deletion",
			userID:      1,
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "user not found",
			userID:      999,
			mockError:   sql.ErrNoRows,
			expectError: true,
		},
		{
			name:        "database error",
			userID:      1,
			mockError:   errors.New("foreign key constraint"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockUserRepositoryInterface(t)
			redisClient, mr := setupTestRedis(t)
			defer mr.Close()

			service := NewUserServiceWithInterface(mockRepo, redisClient)

			mockRepo.On("Delete", tt.userID).Return(tt.mockError)

			err := service.Delete(tt.userID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestUserService_CacheInvalidation tests cache invalidation
func TestUserService_CacheInvalidation(t *testing.T) {
	mockRepo := mocks.NewMockUserRepositoryInterface(t)
	redisClient, mr := setupTestRedis(t)
	defer mr.Close()

	service := NewUserServiceWithInterface(mockRepo, redisClient)

	users := []*models.User{
		{ID: 1, Username: "user1", Email: "user1@example.com"},
	}

	// Populate cache with GetAll
	mockRepo.On("GetAll").Return(users, nil).Once()
	_, err := service.GetAll()
	assert.NoError(t, err)

	// Create a new user - should invalidate cache
	newUser := &models.User{
		ID:       2,
		Username: "newuser",
		Email:    "new@example.com",
	}
	createReq := &models.CreateUserRequest{
		Username: "newuser",
		Email:    "new@example.com",
	}

	mockRepo.On("Create", createReq).Return(newUser, nil)
	_, err = service.Create(createReq)
	assert.NoError(t, err)

	// Next GetAll should hit repository (cache was invalidated)
	mockRepo.On("GetAll").Return(users, nil).Once()
	_, err = service.GetAll()
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
