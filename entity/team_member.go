package entity

import "github.com/google/uuid"

type TeamMember struct {
	TeamMemberID uuid.UUID `json:"team_member_id" gorm:"varchar(36);primaryKey"`
	MemberName   string    `json:"member_name" gorm:"varchar(70);not null"`
	TeamID       uuid.UUID `json:"team_id"`
}
