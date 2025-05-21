package middleware

import (
	"errors"
	"itfest-2025/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (m *middleware) OnlyAdmin(c *gin.Context) {
	user, err := m.jwtAuth.GetLoginUser(c)
	if err != nil {
		response.Error(c, http.StatusForbidden, "failed to get login user", err)
		c.Abort()
		return
	}

	if user.RoleID != 1 {
		response.Error(c, http.StatusForbidden, "this endpoint cannot be access", errors.New("user dont have access"))
		c.Abort()
		return
	}
	c.Next()
}
