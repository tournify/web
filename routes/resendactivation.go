package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/tournify/web/models"
	"log"
	"net/http"
)

func (controller Controller) ResendActivation(c *gin.Context) {
	localize := i18n.NewLocalizer(controller.bundle, domainLanguage(c))

	title, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "resend_activation_email_title",
			Other: "Resend Activation Email",
		},
	})
	pd := controller.defaultPageData(c)
	pd.Title = title
	c.HTML(http.StatusOK, "resendactivation.html", pd)
}

func (controller Controller) ResendActivationPost(c *gin.Context) {
	localize := i18n.NewLocalizer(controller.bundle, domainLanguage(c))

	title, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "resend_activation_email_title",
			Other: "Resend Activation Email",
		},
	})
	pd := controller.defaultPageData(c)
	pd.Title = title
	email := c.PostForm("email")
	user := models.User{Email: email}
	res := controller.db.Where(&user).First(&user)
	if res.Error == nil && user.ActivatedAt == nil {
		activationToken := models.Token{
			Type:    models.TokenUserActivation,
			ModelID: int(user.ID),
		}

		res = controller.db.Where(&activationToken).First(&activationToken)
		if res.Error == nil {
			// If the activation token exists we simply send an email
			go controller.sendActivationEmail(activationToken.Value, user.Email)
		} else {
			// If there is no token then we need to generate a new token
			go controller.activationEmailHandler(user.ID, user.Email)
		}
	} else {
		log.Println(res.Error)
	}

	// We always return a positive response here to prevent user enumeration and other attacks
	pd.Messages = append(pd.Messages, Message{
		Type:    "success",
		Content: "A new activation email has been sent if the account exists and is not already activated. Please remember to check your spam inbox in case the email is not showing in your inbox.",
	})
	c.HTML(http.StatusOK, "resendactivation.html", pd)
}
