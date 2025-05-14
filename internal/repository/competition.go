package repository

import (
	"itfest-2025/entity"

	"gorm.io/gorm"
)

type ICompetitionRepository interface {
	GetCompetitionByID(tx *gorm.DB, competitionID int) (*entity.Competition, error)
}

type CompetitionRepository struct {
	db *gorm.DB
}

func NewCompetitionRepository(db *gorm.DB) ICompetitionRepository {
	return &CompetitionRepository{
		db: db,
	}
}

func (c *CompetitionRepository) GetCompetitionByID(tx *gorm.DB, competitionID int) (*entity.Competition, error) {
	var competition *entity.Competition

	err := tx.Where("competition_id = ?", competitionID).First(&competition).Error
	if err != nil {
		return nil, err
	}

	return competition, nil
}
