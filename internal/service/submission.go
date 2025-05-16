package service

import (
	"errors"
	"itfest-2025/entity"
	"itfest-2025/internal/repository"
	"itfest-2025/model"
	"itfest-2025/pkg/database/mariadb"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ISubmissionService interface {
	GetCurrentStage(userID uuid.UUID) (model.ResStage, error)
	CreateSubmission(userID uuid.UUID, param *model.ReqSubmission) error
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
			IDCurrentStage: 0, 
			NextStage:    firstStage.StageOrder,
			IDNextStage:    firstStage.StageID,
		}
		return data, nil
	} else if err != nil {
		return data, err
	}
	
	// stage, err := s.SubmissionRepository.GetStage(team)
	// if err != nil {
	// 	return data, err
	// }

	nextStage, err := s.SubmissionRepository.GetNextStage(currentStage.StageID, team.CompetitionID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	
	data = model.ResStage{
		IDCurrentStage: currentStage.StageID,
		NextStage: nextStage.StageOrder,
		IDNextStage: nextStage.StageID,
	}

	return data, nil
}

func (s *SubmissionService) CreateSubmission(userID uuid.UUID, param *model.ReqSubmission) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	team, err := s.TeamRepository.GetTeamByUserID(tx, userID)
	if err != nil {
		tx.Rollback()
		return err
	}
	
	newSubmission := &entity.TeamProgress{
		StageID:    param.StageID,    
		Status:     "pending",
		TeamID:     team.TeamID,
		GdriveLink: param.GdriveLink,
	}

	if err := s.SubmissionRepository.CreateSubmission(tx, newSubmission); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}