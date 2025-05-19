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
	UpsertTeam(userID uuid.UUID, param *model.UpsertTeamRequest) (*model.UpsertTeamResponse, error)
	GetMembersByUserID(userID uuid.UUID) (*model.TeamInfoResponse, error)
}

type TeamService struct {
	db                    *gorm.DB
	TeamRepository        repository.ITeamRepository
	CompetitionRepository repository.ICompetitionRepository
}

func NewTeamService(teamRepository repository.ITeamRepository, competitionRepository repository.ICompetitionRepository) ITeamService {
	return &TeamService{
		db:                    mariadb.Connection,
		TeamRepository:        teamRepository,
		CompetitionRepository: competitionRepository,
	}
}

func (t *TeamService) UpsertTeam(userID uuid.UUID, param *model.UpsertTeamRequest) (*model.UpsertTeamResponse, error) {
	if len(param.Members) > 2 {
		return nil, errors.New("maximum of 2 team members allowed")
	}

	tx := t.db.Begin()
	defer tx.Rollback()

	team, err := t.TeamRepository.GetTeamByUserID(tx, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if team == nil {
		teamID := uuid.New()
		newTeam := &entity.Team{
			TeamID:        teamID,
			TeamName:      param.TeamName,
			TeamStatus:    "belum terverifikasi",
			CompetitionID: 1,
			UserID:        userID,
		}

		err := t.TeamRepository.GetTeamByName(tx, param.TeamName)
		if err == nil {
			return nil, errors.New("team name already exists")
		}

		err = t.TeamRepository.CreateTeam(tx, newTeam)
		if err != nil {
			return nil, err
		}
		team = newTeam
	} else {
		team.TeamName = param.TeamName

		err := t.TeamRepository.UpdateTeam(tx, team)
		if err != nil {
			return nil, err
		}

		err = t.TeamRepository.DeleteTeamMembers(tx, team.TeamID)
		if err != nil {
			return nil, err
		}
	}

	for _, v := range param.Members {
		member := &entity.TeamMember{
			TeamMemberID:  uuid.New(),
			TeamID:        team.TeamID,
			MemberName:    v.Name,
			StudentNumber: v.StudentNumber,
		}
		err := t.TeamRepository.CreateTeamMember(tx, member)
		if err != nil {
			return nil, err
		}
	}

	var response model.UpsertTeamResponse
	response.TeamName = team.TeamName
	response.Members = param.Members

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (t *TeamService) GetMembersByUserID(userID uuid.UUID) (*model.TeamInfoResponse, error) {
	tx := t.db.Begin()
	defer tx.Rollback()

	team, err := t.TeamRepository.GetTeamByUserID(tx, userID)
	if err != nil {
		return nil, err
	}

	members, err := t.TeamRepository.GetTeamMemberByTeamID(tx, team.TeamID)
	if err != nil {
		return nil, err
	}

	var memberResponse []model.TeamMembersResponse
	for _, v := range members {
		memberResponse = append(memberResponse, model.TeamMembersResponse{
			FullName:      v.MemberName,
			StudentNumber: v.StudentNumber,
		})
	}

	competition, err := t.CompetitionRepository.GetCompetitionByID(tx, team.CompetitionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model.TeamInfoResponse{
				TeamName: "",
				Members:  []model.TeamMembersResponse{},
			}, nil
		}
		return nil, err
	}

	TeamInforResponse := model.TeamInfoResponse{
		TeamName:            team.TeamName,
		CompetitionCategory: competition.CompetitionName,
		Members:             memberResponse,
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &TeamInforResponse, nil
}
