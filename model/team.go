package model

import "github.com/google/uuid"

type AddTeamMemberRequest struct {
	MemberName string    `json:"member_name" binding:"required"`
	TeamID     uuid.UUID `json:"team_id"`
}

type UpsertTeamRequest struct {
	TeamName string              `json:"team_name" binding:"required"`
	Members  []TeamMemberRequest `json:"members" binding:"required"`
}

type TeamMemberRequest struct {
	Name          string `json:"name" binding:"required"`
	StudentNumber string `json:"student_number" binding:"required"`
}
