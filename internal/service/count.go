package service

import (
	"itfest-2025/internal/repository"
	"itfest-2025/pkg/database/mariadb"

	"gorm.io/gorm"
)

type ICountService interface {
	GetAllCount() (responCount, error)
}

type CountService struct {
	db             *gorm.DB
	TeamRepository repository.ITeamRepository
	UserRepository repository.IUserRepository
}

type responCount struct {
	TotalTeam     int64
	TotalPayment  int64
	TotalBusiness int64
	TotalUIUX     int64
}

func NewCountService(TeamRepository repository.ITeamRepository, UserRepository repository.IUserRepository) *CountService {
	return &CountService{
		db:             mariadb.Connection,
		TeamRepository: TeamRepository,
		UserRepository: UserRepository,
	}
}

func (c *CountService) GetAllCount() (responCount, error) {
	tx := c.db.Begin()
	defer tx.Rollback()

	totalTeam, err := c.TeamRepository.GetCount(tx, "")
	if err != nil {
		return responCount{}, err
	}
	countBusiness, err := c.TeamRepository.GetCount(tx, "3")
	if err != nil {
		return responCount{}, err
	}
	countUIUX, err := c.TeamRepository.GetCount(tx, "2")
	if err != nil {
		return responCount{}, err
	}
	
	countPayment, err := c.UserRepository.GetCountPayment()
	if err != nil {
		return responCount{}, err
	}

	return responCount{
		TotalTeam:     totalTeam,
		TotalPayment:  countPayment,
		TotalBusiness: countBusiness,
		TotalUIUX:     countUIUX,
	}, nil
}
