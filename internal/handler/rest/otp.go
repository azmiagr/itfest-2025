package rest

import (
	"itfest-2025/entity"
	"itfest-2025/model"
	"itfest-2025/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) ResendOtp(c *gin.Context) {
	var req model.GetOtp
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	err = r.service.OtpService.ResendOtp(req)
	if err != nil {
		if err.Error() == "your account is already active" {
			response.Error(c, http.StatusForbidden, "user already verified", err)
			return
		} else if err.Error() == "you can only resend otp every 5 minutes" {
			response.Error(c, http.StatusForbidden, "resend otp failed", err)
			return
		} else {
			response.Error(c, http.StatusInternalServerError, "failed to resend otp", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "success resend otp", nil)
}

func (r *Rest) ResendToken(c *gin.Context) {
	user := c.MustGet("user").(*entity.User)

	err := r.service.OtpService.ResendToken(user.UserID)
	if err != nil {
		if err.Error() == "you can only resend otp every 5 minutes" {
			response.Error(c, http.StatusForbidden, "failed to resend token", err)
			return
		} else {
			response.Error(c, http.StatusInternalServerError, "failed to resend token", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "success resend reset password token", nil)

}
