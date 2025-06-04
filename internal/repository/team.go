package repository

import (
	"itfest-2025/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITeamRepository interface {
	CreateTeam(tx *gorm.DB, team *entity.Team) error
	GetTeamByName(tx *gorm.DB, teamName string) error
	GetTeam(tx *gorm.DB) ([]*entity.Team, error)
	GetTeamByID(tx *gorm.DB, teamID uuid.UUID) (*entity.Team, error)
	CreateTeamMember(tx *gorm.DB, teamMember *entity.TeamMember) error
	GetTeamByUserID(tx *gorm.DB, userID uuid.UUID) (*entity.Team, error)
	UpdateTeam(tx *gorm.DB, team *entity.Team) error
	DeleteTeamMembers(tx *gorm.DB, teamID uuid.UUID) error
	GetTeamMemberByTeamID(tx *gorm.DB, teamID uuid.UUID) ([]*entity.TeamMember, error)
	GetCount(tx *gorm.DB, competitionID string) (int64, error)
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

func (t *TeamRepository) GetCount(tx *gorm.DB, competitionID string) (int64, error) {
	var count int64
	query := tx.Debug().Model(&entity.Team{}).Where("competition_id >= ?", 2)
	if competitionID != "" {
		query = query.Where("competition_id = ?", competitionID)
	}
	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (t *TeamRepository) GetTeam(tx *gorm.DB) ([]*entity.Team, error) {
	var team []*entity.Team
	err := tx.Debug().Preload("Competition").Preload("TeamMember").Find(&team).Error
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (t *TeamRepository) GetTeamByID(tx *gorm.DB, teamID uuid.UUID) (*entity.Team, error) {
	var team entity.Team
	err := tx.Debug().Preload("Competition").Preload("TeamMember").Where("team_id = ?", teamID).First(&team).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (t *TeamRepository) CreateTeamMember(tx *gorm.DB, teamMember *entity.TeamMember) error {
	err := tx.Debug().Create(&teamMember).Error
	if err != nil {
		return err
	}

	return nil
}

func (t *TeamRepository) GetTeamByUserID(tx *gorm.DB, userID uuid.UUID) (*entity.Team, error) {
	var team entity.Team
	err := tx.Where("user_id = ?", userID).First(&team).Error
	if err != nil {
		return nil, err
	}

	return &team, nil
}

func (t *TeamRepository) UpdateTeam(tx *gorm.DB, team *entity.Team) error {
	err := tx.Save(&team).Error
	if err != nil {
		return err
	}

	return nil
}

func (t *TeamRepository) DeleteTeamMembers(tx *gorm.DB, teamID uuid.UUID) error {
	err := tx.Where("team_id = ?", teamID).Delete(&entity.TeamMember{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (t *TeamRepository) GetTeamMemberByTeamID(tx *gorm.DB, teamID uuid.UUID) ([]*entity.TeamMember, error) {
	var members []*entity.TeamMember
	err := tx.Where("team_id = ?", teamID).Find(&members).Error
	if err != nil {
		return nil, err
	}

	return members, nil
}
