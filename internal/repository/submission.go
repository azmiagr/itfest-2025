package repository

import (
	"itfest-2025/entity"

	"gorm.io/gorm"
)

type ISubmissionRepository interface {
	GetFirstStage(competitionID int) (entity.Stages, error)
	GetNextStage(currentOrder int, competitionID int) (entity.Stages, error)
	GetCurrentStage(team *entity.Team) (entity.TeamProgress, error)
	CreateSubmission(tx *gorm.DB, submission *entity.TeamProgress) error
}

type SubmissionRepository struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) ISubmissionRepository {
	return &SubmissionRepository{
		db: db,
	}
}

func (r *SubmissionRepository) GetFirstStage(competitionID int) (entity.Stages, error) {
	var stage entity.Stages
	err := r.db.Where("competition_id = ?", competitionID).
		Order("stage_order ASC").
		First(&stage).Error
	return stage, err
}

func (r *SubmissionRepository) GetNextStage(currentOrder int, competitionID int) (entity.Stages, error) {
	var stage entity.Stages
	err := r.db.Where("competition_id = ? AND stage_order > ?", competitionID, currentOrder).
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
			First(&progress).Error
	err != nil {
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