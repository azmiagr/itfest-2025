package entity

import "github.com/google/uuid"

type Team struct {
	TeamID        uuid.UUID `json:"team_id" gorm:"type:varchar(36);primaryKey"`
	TeamName      string    `json:"team_name" gorm:"type:varchar(50);not null"`
	TeamStatus    string    `json:"team_status" gorm:"type:enum('belum terverifikasi', 'terverifikasi', 'ditolak', 'diproses');not null"`
	UserID        uuid.UUID `json:"user_id"`
	CompetitionID int       `json:"competition_id"`

	TeamMembers    []TeamMember   `json:"team_members" gorm:"foreignKey:TeamID"`
	TeamProgresses []TeamProgress `json:"team_progresses" gorm:"foreignKey:TeamID"`
}
