package model

import "github.com/google/uuid"

type AddTeamMemberRequest struct {
	MemberName string    `json:"member_name" binding:"required"`
	TeamID     uuid.UUID `json:"team_id"`
}

type UpsertTeamRequest struct {
	TeamName string              `json:"team_name"`
	Members  []TeamMemberRequest `json:"members" binding:"required"`
}

type UpsertTeamResponse struct {
	TeamName string              `json:"team_name" binding:"required"`
	Members  []TeamMemberRequest `json:"members" binding:"required"`
}

type TeamMemberRequest struct {
	Name          string `json:"name" binding:"required"`
	StudentNumber string `json:"student_number" binding:"required"`
}

type TeamInfoResponse struct {
	TeamName            string                `json:"team_name"`
	CompetitionCategory string                `json:"competition_category"`
	Members             []TeamMembersResponse `json:"members"`
}

type TeamMembersResponse struct {
	FullName      string `json:"full_name"`
	StudentNumber string `json:"student_number"`
}

type GetAllTeamsResponse struct {
	TeamName        string           `json:"team_name"`
	LeaderName      string           `json:"leader_name"`
	University      string           `json:"university"`
	PaymentStatus   string           `json:"payment_status"`
	CompetitionName string           `json:"competition_name"`
	TeamMembers     []GetTeamMembers `json:"team_members"`
}

type GetTeamMembers struct {
	Name string `json:"name"`
}
