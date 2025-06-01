package rest

import (
	"itfest-2025/pkg/response"
	"net/http"
	"os"

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