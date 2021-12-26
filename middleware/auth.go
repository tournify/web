// Package middleware defines all the middlewares for the application
package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const UserIDKey = "UserID"
const UserRoleKey = "UserRole"

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get(UserIDKey)
		if !exists {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
			return
		}
	}
}

func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get(UserIDKey)
		role, roleExists := c.Get(UserRoleKey)
		if !exists || !roleExists || role != "admin" {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
			return
		}
	}
}
