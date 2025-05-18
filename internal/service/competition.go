package service

import (
	"itfest-2025/internal/repository"
	"itfest-2025/model"
	"itfest-2025/pkg/database/mariadb"

	"gorm.io/gorm"
)

type ICompetitionService interface {
	GetAllCompetitions() ([]*model.GetAllCompetitionsResponse, error)
}

type CompetitionService struct {
	db                    *gorm.DB
	CompetitionRepository repository.ICompetitionRepository
}

func NewCompetitionService(CompetitionRepository repository.ICompetitionRepository) *CompetitionService {
	return &CompetitionService{
		db:                    mariadb.Connection,
		CompetitionRepository: CompetitionRepository,
	}
}

func (c *CompetitionService) GetAllCompetitions() ([]*model.GetAllCompetitionsResponse, error) {
	tx := c.db.Begin()
	defer tx.Rollback()

	competitions, err := c.CompetitionRepository.GetAllCompetitions(tx)
	if err != nil {
		return nil, err
	}

	var response []*model.GetAllCompetitionsResponse
	for _, v := range competitions {
		response = append(response, &model.GetAllCompetitionsResponse{
			CompetitionID:   v.CompetitionID,
			CompetitionName: v.CompetitionName,
			Description:     v.Description,
		})
	}

	return response, nil
}
