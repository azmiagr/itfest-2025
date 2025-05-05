package entity

import (
	"time"

	"github.com/google/uuid"
)

type Announcement struct {
	AnnouncementID uuid.UUID `json:"announcement_id" gorm:"varchar(36);primaryKey"`
	Title          string    `json:"title" gorm:"varchar(255);not null;not null"`
	Description    string    `json:"description" gorm:"text;not null"`
	UserID         uuid.UUID `json:"user_id"`
	CompetitionID  int       `json:"competition_id"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
