package entity

import "github.com/google/uuid"

type Submission struct {
	SubmissionID uuid.UUID `json:"submission_id" gorm:"varchar(36);primaryKey"`
	GdriveLink   string    `json:"gdrive_link" gorm:"varchar(100)"`
	TeamID       uuid.UUID `json:"team_id"`
}
