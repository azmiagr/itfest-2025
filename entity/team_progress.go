package entity

import (
	"time"

	"github.com/google/uuid"
)

type TeamProgress struct {
	TeamProgressID int       `json:"team_progress_id" gorm:"int;primaryKey;autoIncrement"`
	Stage          string    `json:"stage" gorm:"type:enum('tahap 1', 'tahap 2');not null"`
	Status         string    `json:"status" gorm:"type:enum('pending', 'lolos', 'tidak lolos');not null"`
	TeamID         uuid.UUID `json:"team_id"`
	CompetitionID  int       `json:"competition_id"`
	CreatedAt      time.Time `json:"created_at"  gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at"  gorm:"autoUpdateTime"`
}
