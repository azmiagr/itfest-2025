package service

import (
	"itfest-2025/internal/repository"
	"itfest-2025/pkg/database/mariadb"
	"itfest-2025/pkg/template"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type IExcelService interface {
	ExportExcelPayment() (string, error)
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

	// Header style (bold, putih text, background biru, border)
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

	// Style baris ganjil (abu muda) dengan border
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

	// Style baris genap (putih) dengan border
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

	// Style kolom No: hanya alignment tengah + border tanpa fill (agar warna baris tetap muncul)
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
		ColWidths:     map[int]float64{1: 5, 2: 30, 3: 25},
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


	// Kolom No (nomor) style alignment tengah + border, tanpa fill supaya warna baris terlihat
	sheet.ColStyleMap[1] = noColStyle

	fileName, err := template.ExportExcel("Payment IT FEST 2025", []template.ExcelSheet{sheet}, f)
	if err != nil {
		return "", err
	}

	return fileName, nil
}