package model

import "github.com/google/uuid"

type ReqSubmission struct {
	StageID    int       `json:"stage_id" binding:"required"`
	TeamID     uuid.UUID `json:"team_id" binding:"required"`
	GdriveLink string    `json:"gdrive_link" binding:"required,url"`
}

type ResSubmission struct {
	TeamProgressID int       `json:"team_progress_id"`
	Status         string    `json:"status"`
	TeamID         uuid.UUID `json:"team_id"`
	GdriveLink     string    `json:"gdrive_link"`
	CurrentStage   string    `json:"current_stage"`
}

type ResStage struct {
	IDCurrentStage int `json:"id_current_stage"`
	NextStage    int `json:"next_stage"`
	IDNextStage    int `json:"id_next_stage"`
}
