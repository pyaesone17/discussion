package service

import (
	"discussion-forum/internal/models"
	"discussion-forum/internal/service/mocks"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReportService_Create(t *testing.T) {
	tests := []struct {
		name           string
		reportableID   int64
		reportableType string
		request        *models.CreateReportRequest
		mockReturn     *models.Report
		mockError      error
		expectError    bool
	}{
		{
			name:           "successful topic report",
			reportableID:   1,
			reportableType: models.ReportableTypeTopic,
			request: &models.CreateReportRequest{
				UserID: 1,
				Reason: models.ReportReasonSpam,
			},
			mockReturn: &models.Report{
				ID:             1,
				UserID:         1,
				ReportableType: models.ReportableTypeTopic,
				ReportableID:   1,
				Reason:         models.ReportReasonSpam,
				Status:         models.ReportStatusPending,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:           "successful post report",
			reportableID:   5,
			reportableType: models.ReportableTypePost,
			request: &models.CreateReportRequest{
				UserID: 2,
				Reason: models.ReportReasonHarassment,
			},
			mockReturn: &models.Report{
				ID:             2,
				UserID:         2,
				ReportableType: models.ReportableTypePost,
				ReportableID:   5,
				Reason:         models.ReportReasonHarassment,
				Status:         models.ReportStatusPending,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:           "duplicate report error",
			reportableID:   1,
			reportableType: models.ReportableTypeTopic,
			request: &models.CreateReportRequest{
				UserID: 1,
				Reason: models.ReportReasonSpam,
			},
			mockReturn:  nil,
			mockError:   errors.New("Duplicate entry"),
			expectError: true,
		},
		{
			name:           "database error",
			reportableID:   1,
			reportableType: models.ReportableTypeTopic,
			request: &models.CreateReportRequest{
				UserID: 1,
				Reason: models.ReportReasonSpam,
			},
			mockReturn:  nil,
			mockError:   errors.New("database connection error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockReportRepositoryInterface(t)
			service := NewReportServiceWithInterface(mockRepo)

			mockRepo.On("Create", tt.reportableID, tt.reportableType, tt.request).Return(tt.mockReturn, tt.mockError)

			result, err := service.Create(tt.reportableID, tt.reportableType, tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockReturn.ID, result.ID)
				assert.Equal(t, tt.mockReturn.ReportableType, result.ReportableType)
				assert.Equal(t, tt.mockReturn.Reason, result.Reason)
				assert.Equal(t, models.ReportStatusPending, result.Status)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestReportService_GetByID(t *testing.T) {
	tests := []struct {
		name        string
		reportID    int64
		mockReturn  *models.Report
		mockError   error
		expectError bool
	}{
		{
			name:     "successful retrieval",
			reportID: 1,
			mockReturn: &models.Report{
				ID:             1,
				UserID:         1,
				ReportableType: models.ReportableTypeTopic,
				ReportableID:   1,
				Reason:         models.ReportReasonSpam,
				Status:         models.ReportStatusPending,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "report not found",
			reportID:    999,
			mockReturn:  nil,
			mockError:   errors.New("sql: no rows in result set"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockReportRepositoryInterface(t)
			service := NewReportServiceWithInterface(mockRepo)

			mockRepo.On("GetByID", tt.reportID).Return(tt.mockReturn, tt.mockError)

			result, err := service.GetByID(tt.reportID)

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

func TestReportService_GetAll(t *testing.T) {
	tests := []struct {
		name        string
		mockReturn  []*models.Report
		mockError   error
		expectError bool
		expectCount int
	}{
		{
			name: "successful retrieval",
			mockReturn: []*models.Report{
				{
					ID:             1,
					UserID:         1,
					ReportableType: models.ReportableTypeTopic,
					ReportableID:   1,
					Reason:         models.ReportReasonSpam,
					Status:         models.ReportStatusPending,
				},
				{
					ID:             2,
					UserID:         2,
					ReportableType: models.ReportableTypePost,
					ReportableID:   3,
					Reason:         models.ReportReasonHarassment,
					Status:         models.ReportStatusResolved,
				},
			},
			mockError:   nil,
			expectError: false,
			expectCount: 2,
		},
		{
			name:        "empty list",
			mockReturn:  []*models.Report{},
			mockError:   nil,
			expectError: false,
			expectCount: 0,
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
			mockRepo := mocks.NewMockReportRepositoryInterface(t)
			service := NewReportServiceWithInterface(mockRepo)

			mockRepo.On("GetAll").Return(tt.mockReturn, tt.mockError)

			result, err := service.GetAll()

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectCount, len(result))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestReportService_UpdateStatus(t *testing.T) {
	tests := []struct {
		name        string
		reportID    int64
		request     *models.UpdateReportRequest
		mockReturn  *models.Report
		mockError   error
		expectError bool
	}{
		{
			name:     "resolve report",
			reportID: 1,
			request: &models.UpdateReportRequest{
				Status: models.ReportStatusResolved,
			},
			mockReturn: &models.Report{
				ID:             1,
				UserID:         1,
				ReportableType: models.ReportableTypeTopic,
				ReportableID:   1,
				Reason:         models.ReportReasonSpam,
				Status:         models.ReportStatusResolved,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:     "dismiss report",
			reportID: 2,
			request: &models.UpdateReportRequest{
				Status: models.ReportStatusDismissed,
			},
			mockReturn: &models.Report{
				ID:             2,
				UserID:         2,
				ReportableType: models.ReportableTypePost,
				ReportableID:   5,
				Reason:         models.ReportReasonOffTopic,
				Status:         models.ReportStatusDismissed,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:     "report not found",
			reportID: 999,
			request: &models.UpdateReportRequest{
				Status: models.ReportStatusResolved,
			},
			mockReturn:  nil,
			mockError:   errors.New("sql: no rows in result set"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockReportRepositoryInterface(t)
			service := NewReportServiceWithInterface(mockRepo)

			mockRepo.On("UpdateStatus", tt.reportID, tt.request).Return(tt.mockReturn, tt.mockError)

			result, err := service.UpdateStatus(tt.reportID, tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.Status, result.Status)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
