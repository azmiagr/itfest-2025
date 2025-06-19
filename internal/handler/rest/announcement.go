package rest

import (
	"errors"
	"itfest-2025/model"
	"itfest-2025/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) GetAnnouncement(c *gin.Context) {
	data, err := r.service.AnnouncementService.GetAnnouncement()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get announcement", err)
		return
	}

	response.Success(c, http.StatusOK, "success to get announcement", data)	
}

func (r *Rest) CreateAnnouncement(c *gin.Context) {
	var req model.RequestAnnouncement
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}
	
	err = r.service.AnnouncementService.SendAnnouncement(req)
	if err != nil {
		if errors.Is(err, model.ErrUserRecordNotFound) {
			response.Error(c, http.StatusNotFound, "User not found", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to create announcement", err)
		return
	}

	response.Success(c, http.StatusOK, "success to send announcement", nil)
}