package service

import (
	"itfest-2025/internal/repository"
	"itfest-2025/pkg/database/mariadb"
	"itfest-2025/pkg/template"

	"gorm.io/gorm"
)

type IExcelService interface {
	ExportExcelTeam() (string, error)
}

type ExcelService struct {
	db                    *gorm.DB
	TeamRepository        repository.ITeamRepository
	CompetitionRepository repository.ICompetitionRepository
}

func NewExcelService(teamRepository repository.ITeamRepository, competitionRepository repository.ICompetitionRepository) IExcelService {
	return &ExcelService{
		db:                    mariadb.Connection,
		TeamRepository:        teamRepository,
		CompetitionRepository: competitionRepository,
	}
}

func (s *ExcelService) ExportExcelTeam() (string, error) {
	data, err := s.TeamRepository.GetTeam(s.db)
	if err != nil {
		return "", err
	}

	sheet := template.ExcelSheet{
		Name:    "Teams",
		Headers: []string{"Team Name", "Member Name"},
		Rows:    [][]interface{}{},
	}

	for _, team := range data {
		if len(team.TeamMembers) == 0 {
			sheet.Rows = append(sheet.Rows, []interface{}{team.TeamName, ""})
			continue
		}

		for i, member := range team.TeamMembers {
			teamName := ""
			if i == 0 {
				teamName = team.TeamName
			}
			sheet.Rows = append(sheet.Rows, []interface{}{teamName, member.MemberName})
		}
	}

	fileName, err := template.ExportExcel("team_export", []template.ExcelSheet{sheet})
	if err != nil {
		return "", err
	}

	return fileName, nil
}

