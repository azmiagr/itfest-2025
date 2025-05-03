package middleware

import (
	"itfest-2025/model"
	"itfest-2025/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (m *middleware) AuthenticateUser(c *gin.Context) {
	bearer := c.GetHeader("Authorization")
	if bearer == "" {
		response.Error(c, http.StatusUnauthorized, "empty token", nil)
		return
	}

	token := strings.Split(bearer, " ")[1]
	userID, err := m.jwtAuth.ValidateToken(token)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "failed to validate token", err)
		c.Abort()
		return
	}

	user, err := m.service.UserService.GetUser(model.UserParam{
		UserID: userID,
	})
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "failed to get user", err)
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Next()
}
