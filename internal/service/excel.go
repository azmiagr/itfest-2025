package service

import (
	"itfest-2025/internal/repository"
	"itfest-2025/model"
	"itfest-2025/pkg/database/mariadb"
	"itfest-2025/pkg/template"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type IExcelService interface {
	ExportExcelPayment() (string, error)
	ExportExcelTeam() (string, error)
	ExportExcelCompetitionByID(competition int) (string, error)
}

type ExcelService struct {
	db                    *gorm.DB
	UserRepository        repository.IUserRepository
	TeamRepository        repository.ITeamRepository
	CompetitionRepository repository.ICompetitionRepository
}

func NewExcelService(teamRepo repository.ITeamRepository, compRepo repository.ICompetitionRepository, userRepo repository.IUserRepository) IExcelService {
	return &ExcelService{
		db:                    mariadb.Connection,
		TeamRepository:        teamRepo,
		CompetitionRepository: compRepo,
		UserRepository:        userRepo,
	}
}

func (s *ExcelService) ExportExcelPayment() (string, error) {
	data, err := s.UserRepository.GetAllUser()
	if err != nil {
		return "", err
	}

	f := excelize.NewFile()

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4F81BD"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return "", err
	}

	rowStyleOdd, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText: true,
			Vertical: "center",
		},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"F2F2F2"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return "", err
	}

	rowStyleEven, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText: true,
			Vertical: "center",
		},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"FFFFFF"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return "", err
	}

	noColStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return "", err
	}

	headers := []string{"No", "Name", "Email", "NIM", "Team Name", "Link Registrasi", "Payment"}

	sheet := template.ExcelSheet{
		Name:          "Payment",
		Headers:       headers,
		Rows:          [][]interface{}{},
		HeaderStyleID: headerStyle,
		ColWidths:     map[int]float64{1: 5, 2: 30, 3: 25, 4: 30, 5: 30, 6: 100, 7: 100},
		RowStyleMap:   map[int]int{},
		ColStyleMap:   map[int]int{},
	}

	no := 1
	rowIndex := 0
	userColorToggle := 0

	for _, dt := range data {
		sheet.Rows = append(sheet.Rows, []interface{}{no, dt.FullName, dt.Email, dt.StudentNumber, dt.Team.TeamName, dt.RegistrationLink, dt.PaymentTransc})

		excelRowNum := rowIndex + 2

		if userColorToggle == 0 {
			sheet.RowStyleMap[excelRowNum] = rowStyleEven
		} else {
			sheet.RowStyleMap[excelRowNum] = rowStyleOdd
		}

		no++
		rowIndex++
		userColorToggle = 1 - userColorToggle
	}

	sheet.ColStyleMap[1] = noColStyle

	fileName, err := template.ExportExcel("Payment IT FEST 2025", []template.ExcelSheet{sheet}, f)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (s *ExcelService) ExportExcelTeam() (string, error) {
	data, err := s.UserRepository.GetAllUser()
	if err != nil {
		return "", err
	}

	f := excelize.NewFile()

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4F81BD"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	rowStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{WrapText: true, Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	noColStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return "", err
	}

	headers := []string{"No", "Nama User", "Nama Tim", "Nama Kompetisi", "Member"}
	rows := [][]interface{}{}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	teamMap := make(map[string]bool)
	no := 1

	for _, user := range data {
		team, err := s.TeamRepository.GetTeamByUserID(tx, user.UserID)
		if err != nil || team == nil {
			continue
		}

		if teamMap[team.TeamID.String()] {
			continue
		}
		teamMap[team.TeamID.String()] = true

		members, err := s.TeamRepository.GetTeamMemberByTeamID(tx, team.TeamID)
		if err != nil {
			tx.Rollback()
			return "", err
		}

		competition, err := s.CompetitionRepository.GetCompetitionByID(tx, team.CompetitionID)
		if err != nil {
			tx.Rollback()
			return "", err
		}

		for i, member := range members {
			if i == 0 {
				rows = append(rows, []interface{}{
					no,
					user.FullName,
					team.TeamName,
					competition.CompetitionName,
					member.MemberName,
				})
			} else {
				rows = append(rows, []interface{}{
					"", "", "", "", member.MemberName,
				})
			}
		}
		no++
	}

	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	rowStyleMap := make(map[int]int)
	for i := 2; i <= len(rows)+1; i++ {
		rowStyleMap[i] = rowStyle
	}

	colStyleMap := map[int]int{
		1: noColStyle,
	}

	sheet := template.ExcelSheet{
		Name:          "Team",
		Headers:       headers,
		Rows:          rows,
		ColWidths:     map[int]float64{1: 5, 2: 30, 3: 30, 4: 30, 5: 25},
		HeaderStyleID: headerStyle,
		RowStyleMap:   rowStyleMap,
		ColStyleMap:   colStyleMap,
	}

	fileName, err := template.ExportExcel("TeamList IT FEST 2025", []template.ExcelSheet{sheet}, f)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (s *ExcelService) ExportExcelCompetitionByID(competitionID int) (string, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	competition, err := s.CompetitionRepository.GetCompetitionByID(tx, competitionID)
	if err != nil {
		return "", err
	}

	f := excelize.NewFile()

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4F81BD"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	rowStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{WrapText: true, Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	noColStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return "", err
	}

	headers := []string{"No", "Nama User", "Nama Tim", "Nama Kompetisi", "Member"}
	rows := [][]interface{}{}

	no := 1

	for _, team := range competition.Teams {
		user, err := s.UserRepository.GetUser(model.UserParam{
			UserID: team.UserID,
		})
		if err != nil {
			tx.Rollback()
			return "", err
		}

		for i, member := range team.TeamMembers {
			if i == 0 {
				rows = append(rows, []interface{}{
					no,
					user.FullName,
					team.TeamName,
					competition.CompetitionName,
					member.MemberName,
				})
			} else {
				rows = append(rows, []interface{}{
					"", "", "", "", member.MemberName,
				})
			}
		}
		no++
	}

	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	rowStyleMap := make(map[int]int)
	for i := 2; i <= len(rows)+1; i++ {
		rowStyleMap[i] = rowStyle
	}

	colStyleMap := map[int]int{
		1: noColStyle,
	}

	sheet := template.ExcelSheet{
		Name:          "Team",
		Headers:       headers,
		Rows:          rows,
		ColWidths:     map[int]float64{1: 5, 2: 30, 3: 30, 4: 30, 5: 25},
		HeaderStyleID: headerStyle,
		RowStyleMap:   rowStyleMap,
		ColStyleMap:   colStyleMap,
	}

	fileName, err := template.ExportExcel("TeamList IT FEST 2025", []template.ExcelSheet{sheet}, f)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
