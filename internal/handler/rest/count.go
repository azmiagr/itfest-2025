package rest

import (
	"itfest-2025/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) GetCount(c *gin.Context) {
	count, err := r.service.CountService.GetAllCount()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to export team", err)
		return
	}

	response.Success(c, http.StatusOK, "success get count", count)
}