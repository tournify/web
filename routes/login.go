package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/tournify/web/middleware"
	"github.com/tournify/web/models"
	"github.com/tournify/web/util"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

func (controller Controller) Login(c *gin.Context) {
	localize := i18n.NewLocalizer(controller.bundle, domainLanguage(c))

	title, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "login_title",
			Other: "Login",
		},
	})
	pd := controller.defaultPageData(c)
	pd.Title = title
	c.HTML(http.StatusOK, "login.html", pd)
}

func (controller Controller) LoginPost(c *gin.Context) {
	loginError := "Could not login, please make sure that you have typed in the correct email and password. If you have forgotten your password, please click the forgot password link below."
	localize := i18n.NewLocalizer(controller.bundle, domainLanguage(c))

	title, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "login_title",
			Other: "Login",
		},
	})
	pd := controller.defaultPageData(c)
	pd.Title = title
	email := c.PostForm("email")
	user := models.User{Email: email}

	res := controller.db.Where(&user).First(&user)
	if res.Error != nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: loginError,
		})
		log.Println(res.Error)
		c.HTML(http.StatusInternalServerError, "login.html", pd)
		return
	}

	if res.RowsAffected == 0 {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: loginError,
		})
		c.HTML(http.StatusBadRequest, "login.html", pd)
		return
	}

	if user.ActivatedAt == nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: "Account is not activated yet.",
		})
		c.HTML(http.StatusBadRequest, "login.html", pd)
		return
	}

	password := c.PostForm("password")
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: loginError,
		})
		c.HTML(http.StatusBadRequest, "login.html", pd)
		return
	}

	// Generate a ulid for the current session
	sessionIdentifier := util.GenerateULID()

	ses := models.Session{
		Identifier: sessionIdentifier,
	}

	// Session is valid for 1 hour
	ses.ExpiresAt = time.Now().Add(time.Hour)
	ses.UserID = &user.ID

	res = controller.db.Save(&ses)
	if res.Error != nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: loginError,
		})
		log.Println(res.Error)
		c.HTML(http.StatusInternalServerError, "login.html", pd)
		return
	}

	session := sessions.Default(c)
	session.Set(middleware.SessionIdentifierKey, sessionIdentifier)

	err = session.Save()
	if err != nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: loginError,
		})
		log.Println(err)
		c.HTML(http.StatusInternalServerError, "login.html", pd)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/admin")
}
