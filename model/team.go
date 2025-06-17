package model

import (
	"time"

	"github.com/google/uuid"
)

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
	TeamID          string           `json:"team_id"`
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

type ReqUpdateStatusTeam struct {
	TeamID        string `json:"team_id"`
	PaymentStatus string `json:"payment_status" binding:"oneof='belum terverifikasi' 'terverifikasi'"`
}

type TeamInfoResponseAdmin struct {
	TeamName            string                `json:"team_name"`
	CompetitionCategory string                `json:"competition_category"`
	LeaderName          string                `json:"leader_name"`
	StudentNumber       string                `json:"student_number"`
	PaymentStatus       string                `json:"payment_status"`
	PaymentTransc       string                `json:"payment_transaction"`
	Members             []TeamMembersResponse `json:"members"`
	StageNow            StageNow              `json:"progress"`
}

type StageNow struct {
	Stage    string    `json:"stage_name"`
	Deadline time.Time `json:"deadline"`
}

type TeamDetailProgress struct {
	PaymentStatus string   `json:payment_status`
	CurrentStageID int      `json:"current_stageID"`
	CurrentStage  string   `json:"current_stage"`
	NextStage     string   `json:"next_stage"`
	Stages        []Stages `stages`
}

type Stages struct {
	Stage      string    `json:"stage_name"`
	Deadline   time.Time `json:"stage_deadline"`
	GdriveLink string    `json:"link_submission"`
}
