package rest

import (
	"itfest-2025/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) GetAllCompetitions(c *gin.Context) {
	competition, err := r.service.CompetitionService.GetAllCompetitions()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get competitions", err)
		return
	}

	response.Success(c, http.StatusOK, "success to get all competitions", competition)
}
