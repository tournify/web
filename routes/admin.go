package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"net/http"
)

func (controller Controller) Admin(c *gin.Context) {
	localize := i18n.NewLocalizer(controller.bundle, domainLanguage(c))

	title, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "admin_title",
			Other: "Admin",
		},
	})
	pd := controller.defaultPageData(c)
	pd.Title = title
	c.HTML(http.StatusOK, "admin.html", pd)
}
