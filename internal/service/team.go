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

type ITeamService interface {
	AddTeamMember(param model.AddTeamMemberRequest) error
}

type TeamService struct {
	db             *gorm.DB
	TeamRepository repository.ITeamRepository
}

func NewTeamService(teamRepository repository.ITeamRepository) ITeamService {
	return &TeamService{
		db:             mariadb.Connection,
		TeamRepository: teamRepository,
	}
}

func (t *TeamService) AddTeamMember(param model.AddTeamMemberRequest) error {
	tx := t.db.Begin()
	defer tx.Rollback()

	_, err := t.TeamRepository.GetTeamByID(tx, param.TeamID)
	if err != nil {
		return err
	}

	count, err := t.TeamRepository.CountTeamMember(tx, param.TeamID)
	if err != nil {
		return err
	}

	if count >= 2 {
		return errors.New("max member reached")
	}

	teamMemberID := uuid.New()

	member := &entity.TeamMember{
		TeamMemberID: teamMemberID,
		MemberName:   param.MemberName,
		TeamID:       param.TeamID,
	}

	err = t.TeamRepository.CreateTeamMember(tx, member)
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}
