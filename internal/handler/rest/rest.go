package rest

import (
	"fmt"
	"itfest-2025/internal/service"
	"itfest-2025/pkg/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

type Rest struct {
	router     *gin.Engine
	service    *service.Service
	middleware middleware.Interface
}

func NewRest(service *service.Service, middleware middleware.Interface) *Rest {
	return &Rest{
		router:     gin.Default(),
		service:    service,
		middleware: middleware,
	}
}

func (r *Rest) MountEndpoint() {
	r.router.Use(r.middleware.Cors())
	r.router.Use(r.middleware.Timeout())

	routerGroup := r.router.Group("api/v1")
	auth := routerGroup.Group("/auth")

	auth.POST("/register", r.Register)
	auth.PATCH("/register", r.VerifyUser)
	auth.PATCH("/register/resend", r.ResendOtp)
	auth.POST("/login", r.Login)

	user := routerGroup.Group("/users")
	user.Use(r.middleware.AuthenticateUser)
	user.GET("/profile", r.GetUserProfile)
	user.GET("/my-team-info", r.GetTeamInfo)
	user.GET("/my-team-profile", r.GetMyTeamProfile)
	user.POST("/upload-payment", r.UploadPayment)
	user.POST("/forgot-password", r.ForgotPassword)
	user.POST("/verify-token", r.VerifyToken)
	user.PATCH("/update-profile", r.UpdateProfile)
	user.PATCH("/upsert-team", r.UpsertTeam)
	user.PATCH("/change-password", r.ChangePasswordAfterVerify)

	submission := routerGroup.Group("/submissions")
	submission.Use(r.middleware.AuthenticateUser)
	submission.GET("/", r.GetSubmission)
	submission.GET("/stage", r.GetCurrentStage)
	submission.POST("/", r.CreateSubmission)

	competition := routerGroup.Group("/competitions")
	competition.Use(r.middleware.AuthenticateUser)
	competition.POST("/register", r.CompetitionRegistration)
}

func (r *Rest) Run() {
	addr := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")

	r.router.Run(fmt.Sprintf("%s:%s", addr, port))
}
