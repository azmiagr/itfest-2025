package rest

import (
	"itfest-2025/model"
	"itfest-2025/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		response.Error(c, http.StatusInternalServerError, "failed to register user", err)
		return
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
		response.Error(c, http.StatusInternalServerError, "failed to login user", err)
		return
	}

	response.Success(c, http.StatusOK, "success to login user", result)
}

func (r *Rest) UploadPayment(c *gin.Context) {
	userID := c.Param("userID")

	parseID, err := uuid.Parse(userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user ID", err)
		return
	}

	paymentFile, err := c.FormFile("payment")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "payment proof is required", err)
		return
	}

	publicURL, err := r.service.UserService.UploadPayment(parseID, paymentFile)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to upload payment", err)
		return
	}

	response.Success(c, http.StatusOK, "success to upload payment", publicURL)

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
