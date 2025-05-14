package service

import (
	"itfest-2025/internal/repository"
	"itfest-2025/pkg/bcrypt"
	"itfest-2025/pkg/jwt"
	"itfest-2025/pkg/supabase"
)

type Service struct {
	UserService IUserService
	TeamService ITeamService
	OtpService  IOtpService
}

func NewService(repository *repository.Repository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface, supabase supabase.Interface) *Service {
	return &Service{
		UserService: NewUserService(repository.UserRepository, repository.TeamRepository, repository.OtpRepository, bcrypt, jwtAuth, supabase),
		TeamService: NewTeamService(repository.TeamRepository, repository.CompetitionRepository),
		OtpService:  NewOtpService(repository.OtpRepository, repository.UserRepository),
	}
}
