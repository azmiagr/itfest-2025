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
	routerGroup.GET("/competitions", r.GetAllCompetitions)

	auth := routerGroup.Group("/auth")
	auth.POST("/register", r.Register)
	auth.PATCH("/register", r.VerifyUser)
	auth.PATCH("/register/resend", r.ResendOtp)
	auth.POST("/login", r.Login)
	auth.POST("/forgot-password", r.ChangePassword)
	auth.POST("/verify-otp", r.VerifyOtpChangePassword)
	auth.POST("/reset-password", r.ChangePasswordAfterVerify)
	auth.PATCH("/resend-token", r.ResendOtpChangePassword)

	user := routerGroup.Group("/users")
	user.Use(r.middleware.AuthenticateUser)
	user.GET("/profile", r.GetUserProfile)
	user.GET("/my-team-info", r.GetTeamInfo)
	user.GET("/my-team-profile", r.GetMyTeamProfile)
	user.GET("/progress", r.GetProgressByUserID)
	user.POST("/upload-payment", r.UploadPayment)
	user.POST("/change-password", r.ChangePassword)
	user.POST("/verify-token", r.VerifyOtpChangePassword)
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
	competition.POST("/upload-ktm", r.UploadKTM)
	competition.POST("/register/:competition_id", r.CompetitionRegistration)

	admin := routerGroup.Group("/admin")
	admin.Use(r.middleware.AuthenticateUser, r.middleware.OnlyAdmin)
	admin.GET("/payment-status", r.GetUserPaymentStatus)
	admin.GET("/total-participants", r.GetTotalParticipant)
	admin.GET("/count", r.GetCount)
	admin.GET("/teams", r.GetAllTeam)
	admin.GET("/teams/:team_id", r.GetTeamByID)
	admin.GET("/teams/:team_id/progress", r.GetTeamByIDProgress)
	admin.PATCH("/teams/:team_id/progress/:stage_id", r.UpdateStatusSubmission)
	admin.PATCH("/teams/:team_id", r.UpdateTeamStatus)

	announcement := admin.Group("/announcement")
	announcement.GET("/", r.GetAnnouncement)
	announcement.POST("/", r.CreateAnnouncement)

	excel := admin.Group("/excel")
	excel.GET("/data-payment", r.GetExportPayment)
	excel.GET("/data-team", r.GetExportTeam)
	excel.GET("/data-competition", r.GetExportCompetitionID)
}

func (r *Rest) Run() {
	addr := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")

	r.router.Run(fmt.Sprintf("%s:%s", addr, port))
}
