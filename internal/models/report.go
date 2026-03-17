package models

import "time"

const (
	ReportableTypeTopic = "topic"
	ReportableTypePost  = "post"

	ReportReasonSpam           = "spam"
	ReportReasonHarassment     = "harassment"
	ReportReasonOffTopic       = "off_topic"
	ReportReasonMisinformation = "misinformation"
	ReportReasonOther          = "other"

	ReportStatusPending   = "pending"
	ReportStatusResolved  = "resolved"
	ReportStatusDismissed = "dismissed"
)

var ValidReportReasons = map[string]bool{
	ReportReasonSpam:           true,
	ReportReasonHarassment:     true,
	ReportReasonOffTopic:       true,
	ReportReasonMisinformation: true,
	ReportReasonOther:          true,
}

var ValidReportStatuses = map[string]bool{
	ReportStatusPending:   true,
	ReportStatusResolved:  true,
	ReportStatusDismissed: true,
}

type Report struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	ReportableType string    `json:"reportable_type"`
	ReportableID   int64     `json:"reportable_id"`
	Reason         string    `json:"reason"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateReportRequest struct {
	UserID int64  `json:"user_id" binding:"required"`
	Reason string `json:"reason" binding:"required"`
}

type UpdateReportRequest struct {
	Status string `json:"status" binding:"required"`
}
