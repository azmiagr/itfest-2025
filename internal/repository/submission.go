package repository

import (
	"itfest-2025/entity"
	"itfest-2025/model"

	"gorm.io/gorm"
)

type ISubmissionRepository interface {
	GetSubmission(req *model.ReqFilterSubmission) ([]entity.TeamProgress, error)
	GetFirstStage(competitionID int) (entity.Stages, error)
	GetNextStage(currentOrder int, competitionID int) (entity.Stages, error)
	GetCurrentStage(team *entity.Team) (entity.TeamProgress, error)
	CreateSubmission(tx *gorm.DB, submission *entity.TeamProgress) error
	GetStage(tx *gorm.DB, currentID int) (entity.Stages, error)
}

type SubmissionRepository struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) ISubmissionRepository {
	return &SubmissionRepository{
		db: db,
	}
}

func (r *SubmissionRepository) GetSubmission(req *model.ReqFilterSubmission) ([]entity.TeamProgress, error) {
	var dataSubmission []entity.TeamProgress
	query := r.db
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.StageID != 0 {
		query = query.Where("stage_id = ?", req.StageID)
	}
	if req.TeamID != "" {
		query = query.Where("team_id = ?", req.TeamID)
	}
	err := query.Find(&dataSubmission).Error
	return dataSubmission, err
}

func (r *SubmissionRepository) GetFirstStage(competitionID int) (entity.Stages, error) {
	var stage entity.Stages
	err := r.db.Where("competition_id = ?", competitionID).
		Order("stage_order ASC").
		First(&stage).Error
	return stage, err
}

func (r *SubmissionRepository) GetNextStage(currentID int, competitionID int) (entity.Stages, error) {
	var stage entity.Stages
	var currentStage entity.Stages
	err := r.db.First(&currentStage, currentID).Error
	if err != nil {
		return entity.Stages{}, err
	}

	err = r.db.Where("competition_id = ? AND stage_order > ?", competitionID, currentStage.StageOrder).
		Order("stage_order ASC").
		First(&stage).Error
	return stage, err
}

func (t *SubmissionRepository) GetCurrentStage(team *entity.Team) (entity.TeamProgress, error) {
	var progress entity.TeamProgress

	if err := t.db.
		Joins("JOIN stages ON stages.stage_id = team_progresses.stage_id").
		Where("team_id = ?", team.TeamID).
		Order("stages.stage_order DESC").
		First(&progress).Error; err != nil {
		return progress, err
	}

	return progress, nil
}

func (t *SubmissionRepository) CreateSubmission(tx *gorm.DB, submission *entity.TeamProgress) error {
	err := tx.Debug().Create(&submission).Error
	if err != nil {
		return err
	}

	return nil
}

func (t *SubmissionRepository) GetStage(tx *gorm.DB, StageID int) (entity.Stages, error) {
	var stage entity.Stages
	err := tx.First(&stage, StageID).Error
	if err != nil {
		return entity.Stages{}, err
	}

	return stage, nil
}
