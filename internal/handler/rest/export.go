package rest

import (
	"itfest-2025/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) GetExportTeam(c *gin.Context) {
	fileName, err := r.service.ExcelService.ExportExcelTeam()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to export team", err)
		return
	}

	filePath := "public/" + fileName
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.File(filePath)

	response.Success(c, http.StatusCreated, "success to export excel", nil)
}