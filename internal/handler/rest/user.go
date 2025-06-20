package rest

import (
	"itfest-2025/entity"
	"itfest-2025/model"
	"itfest-2025/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *Rest) Register(c *gin.Context) {
	param := model.UserRegister{}
	err := c.ShouldBind(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	token, err := r.service.UserService.Register(&param)
	if err != nil {
		if err.Error() == "email already registered" {
			response.Error(c, http.StatusBadRequest, "failed to register new user", err)
			return
		} else if err.Error() == "password doesn't match" {
			response.Error(c, http.StatusBadRequest, "failed to register new user", err)
			return
		}

	}

	response.Success(c, http.StatusCreated, "success to register new user", token)

}

func (r *Rest) Login(c *gin.Context) {
	param := model.UserLogin{}

	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	result, err := r.service.UserService.Login(param)
	if err != nil {
		if err.Error() == "email or password is wrong" {
			response.Error(c, http.StatusUnauthorized, "email or password is wrong", err)
			return
		} else {
			response.Error(c, http.StatusInternalServerError, "failed to login user", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "success to login user", result)
}

func (r *Rest) UploadPayment(c *gin.Context) {
	user := c.MustGet("user").(*entity.User)

	paymentFile, err := c.FormFile("payment")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "payment proof is required", err)
		return
	}

	publicURL, err := r.service.UserService.UploadPayment(user.UserID, paymentFile)
	if err != nil {
		if err.Error() == "file size exceeds maximum limit of 1MB" {
			response.Error(c, http.StatusBadRequest, "please reduce the file size", err)
		} else {
			response.Error(c, http.StatusInternalServerError, "failed to upload payment", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "success to upload payment", publicURL)

}

func (r *Rest) UploadKTM(c *gin.Context) {
	user := c.MustGet("user").(*entity.User)

	ktmFile, err := c.FormFile("ktm")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ktm is required", err)
		return
	}

	err = r.service.UserService.UploadKTM(user.UserID, ktmFile)
	if err != nil {
		if err.Error() == "file size exceeds maximum limit of 1MB" {
			response.Error(c, http.StatusBadRequest, "please reduce the file size", err)
		} else {
			response.Error(c, http.StatusInternalServerError, "failed to upload payment", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "success to upload ktm", nil)
}

func (r *Rest) VerifyUser(c *gin.Context) {
	var param model.VerifyUser
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	err = r.service.UserService.VerifyUser(param)
	if err != nil {
		if err.Error() == "invalid otp code" {
			response.Error(c, http.StatusUnauthorized, "otp code is wrong", err)
			return
		} else if err.Error() == "otp expired" {
			response.Error(c, http.StatusUnauthorized, "otp code is expired", err)
			return
		} else {
			response.Error(c, http.StatusInternalServerError, "failed to verify user", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "success to verify user", nil)

}

func (r *Rest) UpdateProfile(c *gin.Context) {
	user := c.MustGet("user").(*entity.User)

	var param model.UpdateProfile
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	res, err := r.service.UserService.UpdateProfile(user.UserID, param)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to update user profile", err)
		return
	}

	response.Success(c, http.StatusOK, "success to update user profile", res)
}

func (r *Rest) GetUserProfile(c *gin.Context) {
	user := c.MustGet("user").(*entity.User)

	userProfile, err := r.service.UserService.GetUserProfile(user.UserID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to user profile", err)
		return
	}
	response.Success(c, http.StatusOK, "success to get user profile", userProfile)

}

func (r *Rest) GetMyTeamProfile(c *gin.Context) {
	user := c.MustGet("user").(*entity.User)

	teamProfile, err := r.service.UserService.GetMyTeamProfile(user.UserID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get team profile", err)
		return
	}

	response.Success(c, http.StatusOK, "success to get my team profile", teamProfile)
}

func (r *Rest) ChangePassword(c *gin.Context) {
	var param model.ForgotPasswordRequest
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	token, err := r.service.UserService.ChangePassword(param.Email)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to send email verification", err)
		return
	}

	response.Success(c, http.StatusOK, "success to send email verification password", token)
}

func (r *Rest) VerifyOtpChangePassword(c *gin.Context) {
	var param model.VerifyToken
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	err = r.service.UserService.VerifyOtpChangePassword(param)

	if err != nil {
		if err.Error() == "invalid token" {
			response.Error(c, http.StatusBadRequest, "token is incorrect", err)
			return
		} else if err.Error() == "token expired" {
			response.Error(c, http.StatusBadRequest, "token is already expired", err)
			return
		} else {
			response.Error(c, http.StatusInternalServerError, "failed to verify token", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "success to verify token", nil)
}

func (r *Rest) ChangePasswordAfterVerify(c *gin.Context) {
	var param model.ResetPasswordRequest
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	err = r.service.UserService.ChangePasswordAfterVerify(param)
	if err != nil {
		if err.Error() == "password mismatch" {
			response.Error(c, http.StatusBadRequest, "please check your password", err)
			return
		} else if err.Error() == "new password cannot be same as old password" {
			response.Error(c, http.StatusBadRequest, "please use another password", err)
			return
		} else {
			response.Error(c, http.StatusInternalServerError, "failed to change user password", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "success to change user password", nil)
}

func (r *Rest) CompetitionRegistration(c *gin.Context) {
	user := c.MustGet("user").(*entity.User)

	competitionID := c.Param("competition_id")
	idInt, err := strconv.Atoi(competitionID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to convert competition id", err)
		return
	}

	var param model.CompetitionRegistrationRequest
	err = c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	err = r.service.UserService.CompetitionRegistration(user.UserID, idInt, param)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to register competition", err)
		return
	}

	response.Success(c, http.StatusOK, "success to register competition", nil)
}

func (r *Rest) GetUserPaymentStatus(c *gin.Context) {
	res, err := r.service.UserService.GetUserPaymentStatus()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get user payment status", err)
		return
	}

	response.Success(c, http.StatusOK, "success to get user payment status", res)
}

func (r *Rest) GetTotalParticipant(c *gin.Context) {
	res, err := r.service.UserService.GetTotalParticipant()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get total participant", err)
		return
	}

	response.Success(c, http.StatusOK, "success to get total participant", res)
}
