package entity

import "github.com/google/uuid"

type Team struct {
	TeamID       uuid.UUID    `json:"team_id" gorm:"type:varchar(36);primaryKey"`
	MemberName   string       `json:"name" gorm:"type:varchar(70)"`
	University   string       `json:"university" gorm:"type:varchar(80)"`
	UserID       uuid.UUID    `json:"user_id"`
	Registration Registration `json:"registration" gorm:"foreignKey:TeamID"`
}
