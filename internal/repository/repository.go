package repository

import "gorm.io/gorm"

type Repository struct {
	UserRepository IUserRepository
	TeamRepository ITeamRepository
	OtpRepository  IOtpRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		UserRepository: NewUserRepository(db),
		TeamRepository: NewTeamRepository(db),
		OtpRepository:  NewOtpRepository(db),
	}
}
