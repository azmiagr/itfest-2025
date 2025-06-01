package template

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExcelSheet struct {
	Name          string
	Headers       []string
	Rows          [][]interface{}
	ColWidths     map[int]float64
	HeaderStyleID int
	RowStyleMap   map[int]int
	ColStyleMap   map[int]int
}

func ExportExcel(fileNamePrefix string, sheets []ExcelSheet, f *excelize.File) (string, error) {
	for i, sheet := range sheets {
		sheetName := sheet.Name

		if i == 0 {
			f.SetSheetName("Sheet1", sheetName)
		} else {
			if _, err := f.NewSheet(sheetName); err != nil {
				return "", fmt.Errorf("failed to create sheet %s: %w", sheetName, err)
			}
		}

		// Set header values
		for colIdx, header := range sheet.Headers {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
			f.SetCellValue(sheetName, cell, header)
			if sheet.HeaderStyleID != 0 {
				f.SetCellStyle(sheetName, cell, cell, sheet.HeaderStyleID)
			}
			// Atur lebar kolom jika diset
			colLetter, _ := excelize.ColumnNumberToName(colIdx + 1)
			width := 20.0
			if w, ok := sheet.ColWidths[colIdx+1]; ok {
				width = w
			}
			f.SetColWidth(sheetName, colLetter, colLetter, width)
		}

		// Set row data
		for rowIdx, row := range sheet.Rows {
			excelRow := rowIdx + 2
			for colIdx, val := range row {
				cell, _ := excelize.CoordinatesToCellName(colIdx+1, excelRow)
				f.SetCellValue(sheetName, cell, val)

				// Cek apakah baris punya style
				if styleID, ok := sheet.RowStyleMap[excelRow]; ok {
					f.SetCellStyle(sheetName, cell, cell, styleID)
				}

				// Cek apakah kolom punya style spesifik (misal No rata tengah)
				if styleID, ok := sheet.ColStyleMap[colIdx+1]; ok {
					f.SetCellStyle(sheetName, cell, cell, styleID)
				}
			}
		}
	}

	// Simpan ke folder public/
	outputDir := "public"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	fileName := fmt.Sprintf("%s_%s.xlsx", fileNamePrefix, time.Now().Format("20060102_150405"))
	filePath := filepath.Join(outputDir, fileName)

	if err := f.SaveAs(filePath); err != nil {
		return "", fmt.Errorf("failed to save excel file: %w", err)
	}

	return fileName, nil
}
