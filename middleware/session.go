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
			// TODO if a user hasa large amount of tournaments we may need to optimize this later and not preload every tournament
			res := db.Where(&ses).Preload("User.Tournaments").Preload("Tournaments").First(&ses)
			if res.Error == nil && !ses.HasExpired() {
				c.Set(SessionIDKey, ses.ID)
				if ses.User != nil {
					c.Set(UserIDKey, ses.User.ID)
					c.Set(UserRoleKey, ses.User.Role.Label)
					c.Set(UserTournamentsKey, ses.User.Tournaments)
				} else {
					c.Set(SessionTournamentsKey, ses.Tournaments)
				}
			} else {
				log.Println(res.Error)
			}
		}
		c.Next()
	}
}
