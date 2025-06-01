package rest

import (
	"itfest-2025/pkg/response"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *Rest) GetExportPayment(c *gin.Context) {
	fileName, err := r.service.ExcelService.ExportExcelPayment()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to export team", err)
		return
	}

	filePath := "public/" + fileName
	defer func() {
		_ = os.Remove(filePath)
	}()
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.File(filePath)
}

func (r *Rest) GetExportTeam(c *gin.Context) {
	fileName, err := r.service.ExcelService.ExportExcelTeam()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to export team", err)
		return
	}

	filePath := "public/" + fileName
	defer func() {
		_ = os.Remove(filePath)
	}()
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.File(filePath)
}

func (r *Rest) GetExportCompetitionID(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.Error(c, http.StatusBadRequest, "missing team ID", nil)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid team ID format", err)
		return
	}

	fileName, err := r.service.ExcelService.ExportExcelCompetitionByID(id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to export team", err)
		return
	}

	filePath := "public/" + fileName
	defer func() {
		_ = os.Remove(filePath)
	}()
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.File(filePath)
}