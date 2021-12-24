package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tournify/web/middleware"
	"log"
	"net/http"
)

func (controller Controller) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(middleware.SessionIdentifierKey)
	err := session.Save()
	log.Println(err)

	c.Redirect(http.StatusTemporaryRedirect, "/")
}
