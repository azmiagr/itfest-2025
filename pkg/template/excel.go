package template

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExcelSheet struct {
	Name    string
	Headers []string
	Rows    [][]interface{}
}

func ExportExcel(fileNamePrefix string, sheets []ExcelSheet) (string, error) {
	f := excelize.NewFile()

	for i, sheet := range sheets {
		sheetName := sheet.Name

		if i == 0 {
			f.SetSheetName("Sheet1", sheetName)
		} else {
			if _, err := f.NewSheet(sheetName); err != nil {
				return "", fmt.Errorf("failed to create sheet %s: %w", sheetName, err)
			}
		}

		// Header
		for colIdx, header := range sheet.Headers {
			cell, err := excelize.CoordinatesToCellName(colIdx+1, 1)
			if err != nil {
				return "", fmt.Errorf("failed to get cell name for header: %w", err)
			}
			if err := f.SetCellValue(sheetName, cell, header); err != nil {
				return "", fmt.Errorf("failed to set header value at %s: %w", cell, err)
			}
		}

		// Rows
		for rowIdx, row := range sheet.Rows {
			for colIdx, val := range row {
				cell, err := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
				if err != nil {
					return "", fmt.Errorf("failed to get cell name for data: %w", err)
				}
				if err := f.SetCellValue(sheetName, cell, val); err != nil {
					return "", fmt.Errorf("failed to set cell value at %s: %w", cell, err)
				}
			}
		}
	}

	// Simpan file
	fileName := fmt.Sprintf("%s_%s.xlsx", fileNamePrefix, time.Now().Format("20060102_150405"))
	if err := f.SaveAs("public/"+fileName); err != nil {
		return "", fmt.Errorf("failed to save excel file: %w", err)
	}

	return fileName, nil
}
