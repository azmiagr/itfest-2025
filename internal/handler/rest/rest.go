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
	routerGroup := r.router.Group("api/v1")
	routerGroup.Use(r.middleware.Timeout(), r.middleware.Cors())
	auth := routerGroup.Group("/auth")

	auth.POST("/register", r.Register)
	auth.PATCH("/register", r.VerifyUser)
	auth.PATCH("/register/resend", r.ResendOtp)
	auth.POST("/login", r.Login)

	user := routerGroup.Group("/users")
	user.Use(r.middleware.AuthenticateUser)
	user.GET("/profile", r.GetUserProfile)
	user.GET("/my-team-info", r.GetTeamMember)
	user.GET("/my-team-profile", r.GetMyTeamProfile)
	user.POST("/upload-payment", r.UploadPayment)
	user.PATCH("/update-profile", r.UpdateProfile)
	user.PATCH("/upsert-team", r.UpsertTeam)

	submission := routerGroup.Group("/submissions")
	submission.Use(r.middleware.AuthenticateUser)
	submission.GET("/stage", r.GetCurrentStage)
	submission.POST("/", r.CreateSubmission)
}

func (r *Rest) Run() {
	addr := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")

	r.router.Run(fmt.Sprintf("%s:%s", addr, port))
}
