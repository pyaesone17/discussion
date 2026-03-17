package service

import (
	"discussion-forum/internal/models"
	"discussion-forum/internal/repository"
)

type ReportService struct {
	repo ReportRepositoryInterface
}

func NewReportService(repo *repository.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func NewReportServiceWithInterface(repo ReportRepositoryInterface) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) Create(reportableID int64, reportableType string, req *models.CreateReportRequest) (*models.Report, error) {
	return s.repo.Create(reportableID, reportableType, req)
}

func (s *ReportService) GetByID(id int64) (*models.Report, error) {
	return s.repo.GetByID(id)
}

func (s *ReportService) GetAll() ([]*models.Report, error) {
	return s.repo.GetAll()
}

func (s *ReportService) UpdateStatus(id int64, req *models.UpdateReportRequest) (*models.Report, error) {
	return s.repo.UpdateStatus(id, req)
}
