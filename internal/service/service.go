package service

import (
	"itfest-2025/internal/repository"
	"itfest-2025/pkg/bcrypt"
	"itfest-2025/pkg/jwt"
	"itfest-2025/pkg/supabase"
)

type Service struct {
	UserService        IUserService
	TeamService        ITeamService
	OtpService         IOtpService
	CompetitionService ICompetitionService
	SubmissionService  ISubmissionService
	ExcelService       IExcelService
}

func NewService(repository *repository.Repository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface, supabase supabase.Interface) *Service {
	return &Service{
		UserService:        NewUserService(repository.UserRepository, repository.TeamRepository, repository.OtpRepository, repository.CompetitionRepository, bcrypt, jwtAuth, supabase),
		TeamService:        NewTeamService(repository.TeamRepository, repository.CompetitionRepository),
		OtpService:         NewOtpService(repository.OtpRepository, repository.UserRepository),
		SubmissionService:  NewSubmissionService(repository.SubmissionRepository, repository.TeamRepository),
		CompetitionService: NewCompetitionService(repository.CompetitionRepository),
		ExcelService:       NewExcelService(repository.TeamRepository, repository.CompetitionRepository, repository.UserRepository),
	}
}
