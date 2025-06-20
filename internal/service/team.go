package service

import (
	"errors"
	"itfest-2025/entity"
	"itfest-2025/internal/repository"
	"itfest-2025/model"
	"itfest-2025/pkg/database/mariadb"
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
	tx := t.db.Begin()
	defer tx.Rollback()

	team, err := t.TeamRepository.GetTeamByID(tx, teamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || team == nil {
			return &model.TeamDetailProgress{}, nil
		}
		return nil, err
	}

	competition, err := t.CompetitionRepository.GetCompetitionByID(tx, team.CompetitionID)
	if err != nil {
		return nil, err
	}

	var stages []model.Stages
	dataSubmission, err := t.SubmissionRepository.GetSubmissionAllStage(tx, team.TeamID, team.CompetitionID)
	if err != nil {
		return &model.TeamDetailProgress{}, err
	}
	stages = nil
	if len(dataSubmission) > 0 {
		stages = dataSubmission
	}

	var data model.ResStage
	var currentStageName string
	var nextStageName string

	currentStage, err := t.SubmissionRepository.GetCurrentStage(team)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		firstStage, err := t.SubmissionRepository.GetFirstStage(team.CompetitionID)
		if err != nil {
			return nil, err
		}

		// Jika payment sudah disetujui (terverifikasi), maka stage saat ini dianggap proposal
		if team.TeamStatus == "terverifikasi" {
			data = model.ResStage{
				IDCurrentStage:    firstStage.StageID,
				NextStage:         0, // default kosong, akan diisi jika ada stage selanjutnya
				IDNextStage:       0,
				DeadlineNextStage: time.Time{},
			}

			// Ambil next stage setelah proposal jika ada
			nextStage, err := t.SubmissionRepository.GetNextStage(firstStage.StageID, team.CompetitionID)
			if err == nil {
				data.NextStage = nextStage.StageOrder
				data.IDNextStage = nextStage.StageID
				data.DeadlineNextStage = nextStage.Deadline
				nextStageName = nextStage.StageName
			}

			return &model.TeamDetailProgress{
				TeamCompetition: competition.CompetitionName,
				PaymentStatus:   team.TeamStatus,
				CurrentStageID:  firstStage.StageID,
				CurrentStage:    firstStage.StageName,
				NextStage:       nextStageName,
				Stages:          stages,
			}, nil
		}

		// Kalau belum terverifikasi dan tidak ada stage
		return &model.TeamDetailProgress{
			TeamCompetition: competition.CompetitionName,
			PaymentStatus:   team.TeamStatus,
			CurrentStageID:  0,
			CurrentStage:    "",
			NextStage:       "",
			Stages:          stages,
		}, nil
	} else {
		data = model.ResStage{
			IDCurrentStage: currentStage.StageID,
		}

		// Cek apakah sudah submit untuk currentStage
		submission, err := t.SubmissionRepository.GetSubmission(&model.ReqFilterSubmission{
			StageID: currentStage.StageID,
			TeamID:  team.TeamID.String(),
		})
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		// Jika sudah submit dan statusnya "lolos", geser ke stage berikutnya
		if len(submission) > 0 && submission[0].Status == "lolos" {
			// Ambil next stage dari current
			nextStage, err := t.SubmissionRepository.GetNextStage(currentStage.StageID, team.CompetitionID)
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, err
				}
				stage, err := t.SubmissionRepository.GetStage(tx, currentStage.StageID)
				if err != nil {
					return nil, err
				}
				// Tidak ada next stage, berarti final
				currentStageName = stage.StageName
				nextStageName = ""
			} else {
				// Geser currentStage ke nextStage
				currentStageName = nextStage.StageName
				data.IDCurrentStage = nextStage.StageID

				// Coba ambil stage setelah nextStage
				stageAfterNext, err := t.SubmissionRepository.GetNextStage(nextStage.StageID, team.CompetitionID)
				if err == nil {
					data.NextStage = stageAfterNext.StageOrder
					data.IDNextStage = stageAfterNext.StageID
					data.DeadlineNextStage = stageAfterNext.Deadline
					nextStageName = stageAfterNext.StageName
				} else {
					nextStageName = ""
				}
			}
		} else {
			// Belum submit atau belum lolos, current tetap di sini
			stage, err := t.SubmissionRepository.GetStage(tx, currentStage.StageID)
			if err != nil {
				return nil, err
			}
			currentStageName = stage.StageName

			nextStage, err := t.SubmissionRepository.GetNextStage(currentStage.StageID, team.CompetitionID)
			if err == nil {
				data.NextStage = nextStage.StageOrder
				data.IDNextStage = nextStage.StageID
				data.DeadlineNextStage = nextStage.Deadline
				nextStageName = nextStage.StageName
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
		}
	}

	return &model.TeamDetailProgress{
		TeamCompetition: competition.CompetitionName,
		PaymentStatus:   team.TeamStatus,
		CurrentStageID:  currentStage.StageID,
		CurrentStage:    currentStageName,
		NextStage:       nextStageName,
		Stages:          stages,
	}, nil
}

func (t *TeamService) GetProgressByUserID(userID uuid.UUID) (*model.TeamDetailProgress, error) {
	tx := t.db.Begin()
	defer tx.Rollback()

	team, err := t.TeamRepository.GetTeamByUserID(tx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || team == nil {
			return &model.TeamDetailProgress{}, nil
		}
		return nil, err
	}

	competition, err := t.CompetitionRepository.GetCompetitionByID(tx, team.CompetitionID)
	if err != nil {
		return nil, err
	}

	var stages []model.Stages
	dataSubmission, err := t.SubmissionRepository.GetSubmissionAllStage(tx, team.TeamID, team.CompetitionID)
	if err != nil {
		return &model.TeamDetailProgress{}, err
	}
	stages = nil
	if len(dataSubmission) > 0 {
		stages = dataSubmission
	}

	var data model.ResStage
	var currentStageName string
	var nextStageName string

	currentStage, err := t.SubmissionRepository.GetCurrentStage(team)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		firstStage, err := t.SubmissionRepository.GetFirstStage(team.CompetitionID)
		if err != nil {
			return nil, err
		}

		// Jika payment sudah disetujui (terverifikasi), maka stage saat ini dianggap proposal
		if team.TeamStatus == "terverifikasi" {
			data = model.ResStage{
				IDCurrentStage:    firstStage.StageID,
				NextStage:         0, // default kosong, akan diisi jika ada stage selanjutnya
				IDNextStage:       0,
				DeadlineNextStage: time.Time{},
			}

			// Ambil next stage setelah proposal jika ada
			nextStage, err := t.SubmissionRepository.GetNextStage(firstStage.StageID, team.CompetitionID)
			if err == nil {
				data.NextStage = nextStage.StageOrder
				data.IDNextStage = nextStage.StageID
				data.DeadlineNextStage = nextStage.Deadline
				nextStageName = nextStage.StageName
			}

			return &model.TeamDetailProgress{
				TeamCompetition: competition.CompetitionName,
				PaymentStatus:   team.TeamStatus,
				CurrentStageID:  firstStage.StageID,
				CurrentStage:    firstStage.StageName,
				NextStage:       nextStageName,
				Stages:          stages,
			}, nil
		}

		// Kalau belum terverifikasi dan tidak ada stage
		return &model.TeamDetailProgress{
			TeamCompetition: competition.CompetitionName,
			PaymentStatus:   team.TeamStatus,
			CurrentStageID:  0,
			CurrentStage:    "",
			NextStage:       "",
			Stages:          stages,
		}, nil
	} else {
		data = model.ResStage{
			IDCurrentStage: currentStage.StageID,
		}

		// Cek apakah sudah submit untuk currentStage
		submission, err := t.SubmissionRepository.GetSubmission(&model.ReqFilterSubmission{
			StageID: currentStage.StageID,
			TeamID:  team.TeamID.String(),
		})
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		// Jika sudah submit dan statusnya "lolos", geser ke stage berikutnya
		if len(submission) > 0 && submission[0].Status == "lolos" {
			// Ambil next stage dari current
			nextStage, err := t.SubmissionRepository.GetNextStage(currentStage.StageID, team.CompetitionID)
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, err
				}
				stage, err := t.SubmissionRepository.GetStage(tx, currentStage.StageID)
				if err != nil {
					return nil, err
				}
				// Tidak ada next stage, berarti final
				currentStageName = stage.StageName
				nextStageName = ""
			} else {
				// Geser currentStage ke nextStage
				currentStageName = nextStage.StageName
				data.IDCurrentStage = nextStage.StageID

				// Coba ambil stage setelah nextStage
				stageAfterNext, err := t.SubmissionRepository.GetNextStage(nextStage.StageID, team.CompetitionID)
				if err == nil {
					data.NextStage = stageAfterNext.StageOrder
					data.IDNextStage = stageAfterNext.StageID
					data.DeadlineNextStage = stageAfterNext.Deadline
					nextStageName = stageAfterNext.StageName
				} else {
					nextStageName = ""
				}
			}
		} else {
			// Belum submit atau belum lolos, current tetap di sini
			stage, err := t.SubmissionRepository.GetStage(tx, currentStage.StageID)
			if err != nil {
				return nil, err
			}
			currentStageName = stage.StageName

			nextStage, err := t.SubmissionRepository.GetNextStage(currentStage.StageID, team.CompetitionID)
			if err == nil {
				data.NextStage = nextStage.StageOrder
				data.IDNextStage = nextStage.StageID
				data.DeadlineNextStage = nextStage.Deadline
				nextStageName = nextStage.StageName
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
		}
	}

	return &model.TeamDetailProgress{
		TeamCompetition: competition.CompetitionName,
		PaymentStatus:   team.TeamStatus,
		CurrentStageID:  currentStage.StageID,
		CurrentStage:    currentStageName,
		NextStage:       nextStageName,
		Stages:          stages,
	}, nil
}
