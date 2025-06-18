package service

import (
	"errors"
	"itfest-2025/entity"
	"itfest-2025/internal/repository"
	"itfest-2025/model"
	"itfest-2025/pkg/database/mariadb"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ISubmissionService interface {
	GetSubmission(param *model.ReqFilterSubmission) ([]entity.TeamProgress, error)
	GetCurrentStage(userID uuid.UUID) (model.ResStage, error)
	CreateSubmission(userID uuid.UUID, param *model.ReqSubmission) error
	UpdateStatusSubmission(teamID string, stageID string, param *model.RequestUpdateStatusSubmission) error
}

type SubmissionService struct {
	db                   *gorm.DB
	SubmissionRepository repository.ISubmissionRepository
	TeamRepository       repository.ITeamRepository
}

func NewSubmissionService(submissionRepository repository.ISubmissionRepository, teamRepository repository.ITeamRepository) ISubmissionService {
	return &SubmissionService{
		db:                   mariadb.Connection,
		SubmissionRepository: submissionRepository,
		TeamRepository:       teamRepository,
	}
}

func (s *SubmissionService) GetSubmission(param *model.ReqFilterSubmission) ([]entity.TeamProgress, error) {
	return s.SubmissionRepository.GetSubmission(param)
}

func (s *SubmissionService) GetCurrentStage(userID uuid.UUID) (model.ResStage, error) {
	var data model.ResStage
	tx := s.db.Begin()
	defer tx.Rollback()

	team, err := s.TeamRepository.GetTeamByUserID(tx, userID)
	if err != nil {
		return data, err
	}

	currentStage, err := s.SubmissionRepository.GetCurrentStage(team)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		firstStage, err := s.SubmissionRepository.GetFirstStage(team.CompetitionID)
		if err != nil {
			return data, err
		}

		data = model.ResStage{
			IDCurrentStage:    0,
			NextStage:         firstStage.StageOrder,
			IDNextStage:       firstStage.StageID,
			DeadlineNextStage: firstStage.Deadline,
		}
		return data, nil
	} else if err != nil {
		return data, err
	}
	submission, err := s.SubmissionRepository.GetSubmission(&model.ReqFilterSubmission{
		StageID: data.IDCurrentStage,
		TeamID: team.TeamID.String(),
	})

	if submission[0].Status == "diproses" || submission[0].Status == "tidak lolos" {
		return model.ResStage{
			IDCurrentStage:    currentStage.StageID,
			NextStage:         0,
			IDNextStage:       0,
			DeadlineNextStage: time.Time{},
		}, nil
	}

	nextStage, err := s.SubmissionRepository.GetNextStage(currentStage.StageID, team.CompetitionID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}

	data = model.ResStage{
		IDCurrentStage:    currentStage.StageID,
		NextStage:         nextStage.StageOrder,
		IDNextStage:       nextStage.StageID,
		DeadlineNextStage: nextStage.Deadline,
	}

	return data, nil
}

func (s *SubmissionService) CreateSubmission(userID uuid.UUID, param *model.ReqSubmission) error {
	tx := s.db.Begin()
	defer tx.Rollback()

	stage, err := s.GetCurrentStage(userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	team, err := s.TeamRepository.GetTeamByUserID(tx, userID)
	if err != nil {
		tx.Rollback()
		return err
	}
	submission, err := s.SubmissionRepository.GetSubmission(&model.ReqFilterSubmission{
		StageID: stage.IDCurrentStage,
		TeamID: team.TeamID.String(),
	})

	if submission[0].Status == "tidak lolos" {
		return errors.New("submission ditolak karena stage sebelumnya tidak lolos")
	}
	if submission[0].Status == "diproses" {
		return errors.New("submission sedang diproses")
	}

	if time.Now().After(stage.DeadlineNextStage) {
		return errors.New("submission ditolak karena sudah melewati deadline")
	}
	if team.TeamStatus == "belum terverifikasi" {
		return errors.New("submission ditolak karena belum terverifikasi")
	}

	newSubmission := &entity.TeamProgress{
		StageID:    stage.IDNextStage,
		Status:     "diproses",
		TeamID:     team.TeamID,
		GdriveLink: param.GdriveLink,
	}

	if err := s.SubmissionRepository.CreateSubmission(tx, newSubmission); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *SubmissionService) UpdateStatusSubmission(teamID string, stageID string, param *model.RequestUpdateStatusSubmission) error {
	return s.SubmissionRepository.UpdateStatusSubmission(s.db, teamID, stageID, *param)
}
