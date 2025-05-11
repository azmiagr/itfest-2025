package model

import "github.com/google/uuid"

type AddTeamMemberRequest struct {
	MemberName string    `json:"member_name" binding:"required"`
	TeamID     uuid.UUID `json:"team_id"`
}
