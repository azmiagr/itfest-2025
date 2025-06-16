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
	GetAllTeam() ([]*model.GetAllTeamsResponse, error)
	UpdateTeamStatus(id string, req model.ReqUpdateStatusTeam) error
	GetTeamByID(teamID uuid.UUID) (*model.TeamInfoResponseAdmin, error)
}

type TeamService struct {
	db                    *gorm.DB
	UserRepository        repository.IUserRepository
	TeamRepository        repository.ITeamRepository
	CompetitionRepository repository.ICompetitionRepository
	SubmissionRepository  repository.ISubmissionRepository
}

func NewTeamService(userRepository repository.IUserRepository, teamRepository repository.ITeamRepository, competitionRepository repository.ICompetitionRepository, submissionRepository repository.ISubmissionRepository) ITeamService {
	return &TeamService{
		db:                    mariadb.Connection,
		UserRepository:        userRepository,
		TeamRepository:        teamRepository,
		CompetitionRepository: competitionRepository,
		SubmissionRepository:  submissionRepository,
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

func (t *TeamService) GetAllTeam() ([]*model.GetAllTeamsResponse, error) {
	var (
		res []*model.GetAllTeamsResponse
	)

	tx := t.db.Begin()
	defer tx.Rollback()

	user, err := t.UserRepository.GetAllUser()
	if err != nil {
		return nil, err
	}

	for _, v := range user {
		competition, err := t.CompetitionRepository.GetCompetitionByID(tx, v.Team.CompetitionID)
		if err != nil {
			continue
		}

		var teamMembers []model.GetTeamMembers
		for _, x := range v.Team.TeamMembers {
			teamMembers = append(teamMembers, model.GetTeamMembers{
				Name: x.MemberName,
			})
		}

		res = append(res, &model.GetAllTeamsResponse{
			TeamID:          v.Team.TeamID.String(),
			TeamName:        v.Team.TeamName,
			LeaderName:      v.FullName,
			University:      v.University,
			PaymentStatus:   v.Team.TeamStatus,
			CompetitionName: competition.CompetitionName,
			TeamMembers:     teamMembers,
		})
	}

	return res, nil
}

func (t *TeamService) UpdateTeamStatus(id string, req model.ReqUpdateStatusTeam) error {
	req.TeamID = id
	return t.TeamRepository.UpdateTeamStatus(t.db, req)
}

func (t *TeamService) GetTeamByID(teamID uuid.UUID) (*model.TeamInfoResponseAdmin, error) {
	tx := t.db.Begin()
	defer tx.Rollback()

	team, err := t.TeamRepository.GetTeamByID(tx, teamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || team == nil {
			return &model.TeamInfoResponseAdmin{}, nil
		}
		return nil, err
	}

	user, err := t.UserRepository.GetUser(model.UserParam{
		UserID: team.UserID,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || user == nil {
			user = &entity.User{}
		} else {
			return nil, err
		}
	}

	members, err := t.TeamRepository.GetTeamMemberByTeamID(tx, team.TeamID)
	if err != nil {
		members = []*entity.TeamMember{}
	}

	var memberResponse []model.TeamMembersResponse
	for _, v := range members {
		memberResponse = append(memberResponse, model.TeamMembersResponse{
			FullName:      v.MemberName,
			StudentNumber: v.StudentNumber,
		})
	}

	competition, err := t.CompetitionRepository.GetCompetitionByID(tx, team.CompetitionID)
	competitionName := ""
	if err == nil && competition != nil {
		competitionName = competition.CompetitionName
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var data model.ResStage
	currentStage, err := t.SubmissionRepository.GetCurrentStage(team)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		firstStage, err := t.SubmissionRepository.GetFirstStage(team.CompetitionID)
		if err != nil {
			return nil, err
		}

		data = model.ResStage{
			IDCurrentStage:    0,
			NextStage:         firstStage.StageOrder,
			IDNextStage:       firstStage.StageID,
			DeadlineNextStage: firstStage.Deadline,
		}
	} else {
		nextStage, err := t.SubmissionRepository.GetNextStage(currentStage.StageID, team.CompetitionID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	
		data = model.ResStage{
			IDCurrentStage:    currentStage.StageID,
			NextStage:         nextStage.StageOrder,
			IDNextStage:       nextStage.StageID,
			DeadlineNextStage: nextStage.Deadline,
		}
	}

	stage, err := t.SubmissionRepository.GetStage(tx, data.NextStage)
	if err != nil {
		return nil, err
	}

	response := model.TeamInfoResponseAdmin{
		TeamName:            team.TeamName,
		CompetitionCategory: competitionName,
		LeaderName:          user.FullName,
		StudentNumber:       user.StudentNumber,
		PaymentStatus:       team.TeamStatus,
		PaymentTransc:       user.PaymentTransc,
		Members:             memberResponse,
		StageNow: model.StageNow{
			Stage:    stage.StageName,
			Deadline: data.DeadlineNextStage,
		},
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &response, nil
}
