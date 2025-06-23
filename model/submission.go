package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUnverifiedAccount    = errors.New("submission ditolak karena belum terverifikasi atau ditolak")
	ErrNotPassedPrevious    = errors.New("submission ditolak karena stage sebelumnya tidak lolos")
	ErrSubmissionProcessing = errors.New("submission sedang diproses")
	ErrPassedDeadline       = errors.New("submission ditolak karena sudah melewati deadline")
	ErrNoStage              = errors.New("submission ditolak karena stage tidak tersedia")
)

type ReqSubmission struct {
	GdriveLink string `json:"gdrive_link" binding:"required,url"`
}
type ReqFilterSubmission struct {
	StageID int    `form:"stage_id" json:"stage_id"`
	Status  string `form:"status" json:"status"`
	TeamID  string `form:"team_id" json:"team_id"`
}
type ResSubmission struct {
	TeamProgressID int       `json:"team_progress_id"`
	Status         string    `json:"status"`
	TeamID         uuid.UUID `json:"team_id"`
	GdriveLink     string    `json:"gdrive_link"`
	CurrentStage   string    `json:"current_stage"`
}

type ResStage struct {
	IDCurrentStage    int       `json:"id_current_stage"`
	NextStage         int       `json:"next_stage"`
	IDNextStage       int       `json:"id_next_stage"`
	DeadlineNextStage time.Time `json:"deadline_next_stage"`
}

type ResCurrentSubmission struct {
	TeamProgressID int       `json:"team_progress_id"`
	Stage          string    `json:"stage"`
	Deadline       time.Time `json:"deadline"`
	Status         string    `json:"status"`
	GdriveLink     string    `json:"gdrive_link"`
}

type RequestUpdateStatusSubmission struct {
	SubmissionStatus string `json:"submission_status" binding:"oneof='diproses' 'lolos' 'tidak lolos'"`
}
