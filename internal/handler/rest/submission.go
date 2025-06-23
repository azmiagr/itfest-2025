package rest

import (
	"errors"
	"itfest-2025/entity"
	"itfest-2025/model"
	"itfest-2025/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) GetSubmission(c *gin.Context) {
	param := &model.ReqFilterSubmission{}
	err := c.ShouldBindQuery(param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	data, err := r.service.SubmissionService.GetSubmission(param)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get submission", err)
		return
	}

	response.Success(c, http.StatusCreated, "success to get submission", data)
}

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
		if errors.Is(err, model.ErrUnverifiedAccount) {
			response.Error(c, http.StatusForbidden, "Status team ditolak atau belum diverifikasi", err)
			return
		} else if errors.Is(err, model.ErrNotPassedPrevious) {
			response.Error(c, http.StatusUnprocessableEntity, "submission failed", err)
			return
		} else if errors.Is(err, model.ErrSubmissionProcessing) {
			response.Error(c, http.StatusConflict, "submission sedang diproses", err)
			return
		} else if errors.Is(err, model.ErrPassedDeadline) {
			response.Error(c, http.StatusGone, "submission melewati deadline", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to create submission", err)
		return
	}

	response.Success(c, http.StatusCreated, "success to create new submission", nil)
}

func (r *Rest) UpdateStatusSubmission(c *gin.Context) {
	var req model.RequestUpdateStatusSubmission
	teamID := c.Param("team_id")
	stageID := c.Param("stage_id")

	if teamID == "" {
		response.Error(c, http.StatusBadRequest, "team ID is invalid", nil)
		return
	}
	if stageID == "" || stageID == "0" {
		response.Error(c, http.StatusBadRequest, "stage ID is invalid", nil)
		return
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	if req.SubmissionStatus != "lolos" && req.SubmissionStatus != "tidak lolos" {
		response.Error(c, http.StatusBadRequest, "invalid Submission status", nil)
		return
	}

	err = r.service.SubmissionService.UpdateStatusSubmission(teamID, stageID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to update team status", err)
		return
	}

	response.Success(c, http.StatusOK, "success update team status", nil)
}