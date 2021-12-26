package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tournify/web/models"
	"gorm.io/gorm"
	"log"
)

func Session(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionIdentifierInterface := session.Get(SessionIdentifierKey)

		if sessionIdentifier, ok := sessionIdentifierInterface.(string); ok {
			ses := models.Session{
				Identifier: sessionIdentifier,
			}
			res := db.Where(&ses).Preload("User").Preload("Role").First(&ses)
			if res.Error == nil && !ses.HasExpired() {
				c.Set(UserIDKey, ses.User.ID)
				c.Set(UserRoleKey, ses.User.Role.Label)
			} else {
				log.Println(res.Error)
			}
		}
		c.Next()
	}
}
