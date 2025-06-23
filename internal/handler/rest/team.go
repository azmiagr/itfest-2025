package rest

import (
	"errors"
	"itfest-2025/entity"
	"itfest-2025/model"
	"itfest-2025/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *Rest) UpsertTeam(c *gin.Context) {
	var param model.UpsertTeamRequest
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	user := c.MustGet("user").(*entity.User)

	res, err := r.service.TeamService.UpsertTeam(user.UserID, &param)
	if err != nil {
		if err.Error() == "maximum of 2 team members allowed" {
			response.Error(c, http.StatusBadRequest, "cannot add another team member", err)
			return
		} else if err.Error() == "team name already exists" {
			response.Error(c, http.StatusBadRequest, "cannot use this team name", err)
			return
		} else {
			response.Error(c, http.StatusInternalServerError, "failed to upsert team", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "success upsert team", res)
}

func (r *Rest) GetTeamInfo(c *gin.Context) {
	user := c.MustGet("user").(*entity.User)

	teamInfo, err := r.service.TeamService.GetMembersByUserID(user.UserID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get team members", err)
		return
	}

	response.Success(c, http.StatusOK, "success get team members", teamInfo)

}

func (r *Rest) GetAllTeam(c *gin.Context) {
	res, err := r.service.TeamService.GetAllTeam()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get all team informations", err)
		return
	}

	response.Success(c, http.StatusOK, "success get all team informations", res)
}

func (r *Rest) UpdateTeamStatus(c *gin.Context) {
	var req model.ReqUpdateStatusTeam
	teamID := c.Param("team_id")

	if teamID == "" {
		response.Error(c, http.StatusBadRequest, "team ID is required", nil)
		return
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	err = r.service.TeamService.UpdateTeamStatus(teamID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to update team status", err)
		return
	}

	response.Success(c, http.StatusOK, "success update team status", nil)
}

func (r *Rest) GetTeamByID(c *gin.Context) {
	teamIDParam := c.Param("team_id")

	if teamIDParam == "" {
		response.Error(c, http.StatusBadRequest, "team ID is required", nil)
		return
	}
	teamID, _ := uuid.Parse(teamIDParam)

	data, err := r.service.TeamService.GetTeamByID(teamID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get team", err)
		return
	}

	response.Success(c, http.StatusOK, "success get team", data)
}

func (r *Rest) GetTeamByIDProgress(c *gin.Context) {
	teamIDParam := c.Param("team_id")

	if teamIDParam == "" {
		response.Error(c, http.StatusBadRequest, "team ID is required", nil)
		return
	}
	teamID, _ := uuid.Parse(teamIDParam)

	data, err := r.service.TeamService.GetDetailTeam(teamID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get progress team", err)
		return
	}

	response.Success(c, http.StatusOK, "success get progress team", data)
}

func (r *Rest) GetProgressByUserID(c *gin.Context) {
	userID := c.MustGet("user").(*entity.User)

	data, err := r.service.TeamService.GetProgressByUserID(userID.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "failed to get progress team", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to get progress team", err)
		return
	}

	response.Success(c, http.StatusOK, "success get progress team", data)
}