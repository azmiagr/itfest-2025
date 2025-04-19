package entity

import "github.com/google/uuid"

type Registration struct {
	RegistrationID uuid.UUID `json:"registration_id" gorm:"type:varchar(36);primaryKey"`
	CompetitionID  int       `json:"competition_id"`
	TeamID         uuid.UUID `json:"team_id"`
}
