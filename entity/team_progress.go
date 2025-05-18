package entity

import (
	"time"

	"github.com/google/uuid"
)

type TeamProgress struct {
	TeamProgressID int       `json:"team_progress_id" gorm:"int;primaryKey;autoIncrement"`
	StageID        int       `json:"stage_id"`
	Status         string    `json:"status" gorm:"type:enum('diproses', 'lolos', 'tidak lolos');not null"`
	TeamID         uuid.UUID `json:"team_id"`
	GdriveLink     string    `json:"gdrive_link" gorm:"varchar(100);not null"`
	CreatedAt      time.Time `json:"created_at"  gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at"  gorm:"autoUpdateTime"`
}
