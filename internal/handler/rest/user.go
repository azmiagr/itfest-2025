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

	userID, err := r.service.UserService.Register(&param)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to register user", err)
		return
	}

	loginResponse := model.RegisterResponse{
		UserID: userID,
	}

	response.Success(c, http.StatusCreated, "success to register new user", loginResponse)

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
