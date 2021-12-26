package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"net/http"
)

func (controller Controller) TermsOfService(c *gin.Context) {
	localize := i18n.NewLocalizer(controller.bundle, domainLanguage(c))

	title, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "terms_of_service_title",
			Other: "Terms of Service",
		},
	})
	pd := controller.defaultPageData(c)
	pd.Title = title
	c.HTML(http.StatusOK, "terms-of-service.html", pd)
}

func (controller Controller) PrivacyPolicy(c *gin.Context) {
	localize := i18n.NewLocalizer(controller.bundle, domainLanguage(c))

	title, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "privacy_policy_title",
			Other: "Privacy Policy",
		},
	})
	pd := controller.defaultPageData(c)
	pd.Title = title
	c.HTML(http.StatusOK, "privacy-policy.html", pd)
}
