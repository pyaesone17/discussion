package repository

import (
	"database/sql"
	"discussion-forum/internal/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) Create(reportableID int64, reportableType string, req *models.CreateReportRequest) (*models.Report, error) {
	result, err := r.db.Exec(`
		INSERT INTO reports (user_id, reportable_type, reportable_id, reason)
		VALUES (?, ?, ?, ?)
	`, req.UserID, reportableType, reportableID, req.Reason)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *ReportRepository) GetByID(id int64) (*models.Report, error) {
	report := &models.Report{}
	err := r.db.QueryRow(`
		SELECT id, user_id, reportable_type, reportable_id, reason, status, created_at, updated_at
		FROM reports WHERE id = ?
	`, id).Scan(
		&report.ID, &report.UserID, &report.ReportableType, &report.ReportableID,
		&report.Reason, &report.Status, &report.CreatedAt, &report.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return report, nil
}

func (r *ReportRepository) GetAll() ([]*models.Report, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, reportable_type, reportable_id, reason, status, created_at, updated_at
		FROM reports ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []*models.Report
	for rows.Next() {
		report := &models.Report{}
		err := rows.Scan(
			&report.ID, &report.UserID, &report.ReportableType, &report.ReportableID,
			&report.Reason, &report.Status, &report.CreatedAt, &report.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	return reports, nil
}

func (r *ReportRepository) UpdateStatus(id int64, req *models.UpdateReportRequest) (*models.Report, error) {
	_, err := r.db.Exec(`
		UPDATE reports SET status = ? WHERE id = ?
	`, req.Status, id)
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}
