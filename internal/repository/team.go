package repository

import (
	"itfest-2025/entity"

	"gorm.io/gorm"
)

type ITeamRepository interface {
	CreateTeam(tx *gorm.DB, team *entity.Team) error
	GetTeamByName(tx *gorm.DB, teamName string) error
}

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) ITeamRepository {
	return &TeamRepository{
		db: db,
	}
}

func (t *TeamRepository) CreateTeam(tx *gorm.DB, team *entity.Team) error {
	err := tx.Debug().Create(&team).Error
	if err != nil {
		return err
	}

	return nil
}

func (t *TeamRepository) GetTeamByName(tx *gorm.DB, teamName string) error {
	var team entity.Team
	err := tx.Debug().Where("team_name = ?", teamName).First(&team).Error
	if err != nil {
		return err
	}
	return nil
}
