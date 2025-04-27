package rest

import (
	"fmt"
	"itfest-2025/internal/service"
	"os"

	"github.com/gin-gonic/gin"
)

type Rest struct {
	router  *gin.Engine
	service *service.Service
}

func NewRest(service *service.Service) *Rest {
	return &Rest{
		router:  gin.Default(),
		service: service,
	}
}

func (r *Rest) MountEndpoint() {
	routerGroup := r.router.Group("api/v1")
	routerGroup.POST("/register", r.Register)
	routerGroup.POST("/login", r.Login)
	routerGroup.POST("/upload-payment/:userID", r.UploadPayment)
}

func (r *Rest) Run() {
	addr := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")

	r.router.Run(fmt.Sprintf("%s:%s", addr, port))
}
