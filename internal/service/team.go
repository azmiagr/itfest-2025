package service

import (
	"errors"
	"itfest-2025/entity"
	"itfest-2025/internal/repository"
	"itfest-2025/model"
	"itfest-2025/pkg/database/mariadb"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITeamService interface {
	UpsertTeam(userID uuid.UUID, param *model.UpsertTeamRequest) (*model.UpsertTeamResponse, error)
	GetMembersByUserID(userID uuid.UUID) (*model.TeamInfoResponse, error)
	GetAllTeam() ([]*model.GetAllTeamsResponse, error)
	UpdateTeamStatus(id string, req model.ReqUpdateStatusTeam) error
	GetTeamByID(teamID uuid.UUID) (*model.TeamInfoResponseAdmin, error)
	GetDetailTeam(teamID uuid.UUID) (*model.TeamDetailProgress, error)
	GetProgressByUserID(userID uuid.UUID) (*model.TeamDetailProgress, error)
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
		if v.RoleID == 1 {
			continue
		}

		var teamMembers []model.GetTeamMembers
		var dataStage model.TeamDetailProgress
		dataStage.CurrentStage = "Tidak memiliki stage"
		for _, x := range v.Team.TeamMembers {
			teamMembers = append(teamMembers, model.GetTeamMembers{
				Name: x.MemberName,
			})
		}
		if v.Team.CompetitionID > 1 {
			dataCurrent, _ := t.getProgress(v.Team.TeamID, true)
			if err == nil {
				dataStage.CurrentStage = dataCurrent.CurrentStage
			}
		}
		res = append(res, &model.GetAllTeamsResponse{
			TeamID:          v.Team.TeamID.String(),
			TeamName:        v.Team.TeamName,
			LeaderName:      v.FullName,
			University:      v.University,
			PaymentStatus:   v.Team.TeamStatus,
			CompetitionName: competition.CompetitionName,
			CurrentStage:    dataStage.CurrentStage,
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
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &model.TeamInfoResponseAdmin{
					TeamName:            team.TeamName,
					CompetitionCategory: competitionName,
					LeaderName:          user.FullName,
					StudentNumber:       user.StudentNumber,
					PaymentStatus:       team.TeamStatus,
					PaymentTransc:       user.PaymentTransc,
					Members:             memberResponse,
					StageNow:            model.StageNow{},
				}, nil
			}
			return nil, err
		}

		data = model.ResStage{
			IDCurrentStage:    1,
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

	stage, err := t.SubmissionRepository.GetStage(tx, data.IDCurrentStage)
	if err != nil {
		stage = entity.Stages{}
	}

	submission := ""
	dataSubmission, err := t.SubmissionRepository.GetSubmission(&model.ReqFilterSubmission{
		TeamID:  team.TeamID.String(),
		StageID: stage.StageID,
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			submission = "diproses"
		} else {
			return nil, err
		}
	} else if len(dataSubmission) == 0 {
		submission = "diproses"
	} else {
		submission = dataSubmission[0].Status
	}

	if team.TeamStatus != "terverifikasi" {
		submission = "Akun belum terverifikasi"
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
			Status:   submission,
			Deadline: stage.Deadline,
		},
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &response, nil
}

func (t *TeamService) GetDetailTeam(teamID uuid.UUID) (*model.TeamDetailProgress, error) {
	return t.getProgress(teamID, true)
}

func (t *TeamService) GetProgressByUserID(userID uuid.UUID) (*model.TeamDetailProgress, error) {
	return t.getProgress(userID, false)
}

func (t *TeamService) getProgress(ID uuid.UUID, isAdmin bool) (*model.TeamDetailProgress, error) {
	tx := t.db.Begin()
	defer tx.Rollback()

	team, err := t.TeamRepository.GetTeamByUserID(tx, ID)
	if isAdmin {
		team, err = t.TeamRepository.GetTeamByID(tx, ID)
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || team == nil {
			return &model.TeamDetailProgress{}, nil
		}
		return nil, err
	}

	user, err := t.UserRepository.GetUser(model.UserParam{UserID: team.UserID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || user == nil {
			user = &entity.User{}
		} else {
			return nil, err
		}
	}

	competition, err := t.CompetitionRepository.GetCompetitionByID(tx, team.CompetitionID)
	if err != nil {
		return nil, err
	}

	competitionName := strings.ToLower(competition.CompetitionName)
	isBP := strings.Contains(competitionName, "bp") || strings.Contains(competitionName, "business")
	isUIUX := strings.Contains(competitionName, "ui") || strings.Contains(competitionName, "ux")

	stages, err := t.SubmissionRepository.GetSubmissionAllStage(tx, team.TeamID, team.CompetitionID)
	if err != nil {
		return nil, err
	}

	// dummy payment stage
	paymentStage := model.Stages{
		Stage:      "Payment",
		GdriveLink: user.PaymentTransc,
		Status:     team.TeamStatus,
		Deadline:   time.Time{},
	}

	if isBP {
		if len(stages) >= 2 {
			paymentStage.Deadline = stages[1].Deadline
		}
		if len(stages) >= 1 {
			stages = append(stages[:1], append([]model.Stages{paymentStage}, stages[1:]...)...)
		} else {
			stages = append(stages, paymentStage)
		}
	} else if isUIUX {
		if len(stages) >= 1 {
			paymentStage.Deadline = stages[0].Deadline
		}
		stages = append([]model.Stages{paymentStage}, stages...)
	}

	index := 0
	for index < len(stages) {
		stageName := strings.ToLower(stages[index].Stage)
		stageStatus := strings.ToLower(stages[index].Status)

		// Jika stage = payment dan sudah terverifikasi, skip
		if stageName == "payment" && team.TeamStatus == "terverifikasi" {
			index++
			continue
		}

		// Jika belum lolos, berarti current stage
		if stageStatus != "lolos" {
			break
		}

		index++
	}

	if index >= len(stages) {
		index = len(stages) - 1
	}

	currentStage := stages[index]
	nextStage := ""
	if index+1 < len(stages) {
		nextStage = stages[index+1].Stage
	}

	return &model.TeamDetailProgress{
		TeamCompetition: competition.CompetitionName,
		PaymentStatus:   team.TeamStatus,
		CurrentStageID:  currentStage.StageID,
		CurrentStage:    currentStage.Stage,
		NextStage:       nextStage,
		Stages:          stages,
	}, nil
}
