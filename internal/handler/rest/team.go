package rest

import (
	"itfest-2025/entity"
	"itfest-2025/model"
	"itfest-2025/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) AddTeamMember(c *gin.Context) {
	var param model.AddTeamMemberRequest
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	user := c.MustGet("user").(*entity.User)
	param.TeamID = user.Team.TeamID

	err = r.service.TeamService.AddTeamMember(param)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to add team member", err)
		return
	}

	response.Success(c, http.StatusCreated, "success add team member", nil)
}

func (r *Rest) UpsertTeam(c *gin.Context) {
	var param model.UpsertTeamRequest
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	user := c.MustGet("user").(*entity.User)

	err = r.service.TeamService.UpsertTeam(user.UserID, &param)
	if err != nil {
		if err.Error() == "maximum of 2 team members allowed" {
			response.Error(c, http.StatusBadRequest, "cannot add another team member", err)
			return
		} else {
			response.Error(c, http.StatusInternalServerError, "failed to upsert team", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "success upsert team", nil)
}
