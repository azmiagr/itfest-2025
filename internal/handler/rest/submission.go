package rest

import (
	"itfest-2025/entity"
	"itfest-2025/model"
	"itfest-2025/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) GetCurrentStage(c *gin.Context) {
	user := c.MustGet("user").(*entity.User)
	
	data, err := r.service.SubmissionService.GetCurrentStage(user.UserID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get current stage", err)
		return
	}

	response.Success(c, http.StatusCreated, "success to get current stage", data)
}

func (r *Rest) CreateSubmission(c *gin.Context) {
	param := model.ReqSubmission{}
	user := c.MustGet("user").(*entity.User)
	
	err := c.ShouldBind(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	err = r.service.SubmissionService.CreateSubmission(user.UserID, &param)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create submission", err)
		return
	}

	response.Success(c, http.StatusCreated, "success to create new submission", nil)
}