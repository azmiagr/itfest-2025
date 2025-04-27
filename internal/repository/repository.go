package repository

import "gorm.io/gorm"

type Repository struct {
	UserRepository IUserRepository
	TeamRepository ITeamRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		UserRepository: NewUserRepository(db),
		TeamRepository: NewTeamRepository(db),
	}
}
