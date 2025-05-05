package entity

import "github.com/google/uuid"

type Team struct {
	TeamID       uuid.UUID    `json:"team_id" gorm:"type:varchar(36);primaryKey"`
	TeamName     string       `json:"team_name" gorm:"type:varchar(50);not null"`
	Major        string       `json:"major" gorm:"type:varchar(80);not null"`
	University   string       `json:"university" gorm:"type:varchar(80);not null"`
	TeamStatus   string       `json:"team_status" gorm:"type:enum('belum terverifikasi', 'terverifikasi');not null"`
	UserID       uuid.UUID    `json:"user_id"`
	Registration Registration `json:"registration" gorm:"foreignKey:TeamID"`

	TeamMembers    []TeamMember   `json:"team_members" gorm:"foreignKey:TeamID"`
	Submissions    []Submission   `json:"submissions" gorm:"foreignKey:TeamID"`
	TeamProgresses []TeamProgress `json:"team_progresses" gorm:"foreignKey:TeamID"`
}
