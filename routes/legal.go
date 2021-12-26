package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (controller Controller) TermsOfService(c *gin.Context) {
	pd := PageData{
		Title:           "Terms of Service",
		IsAuthenticated: isAuthenticated(c),
		IsAdmin:         isAdmin(c),
		CacheParameter:  controller.config.CacheParameter,
	}
	c.HTML(http.StatusOK, "terms-of-service.html", pd)
}

func (controller Controller) PrivacyPolicy(c *gin.Context) {
	pd := PageData{
		Title:           "Privacy Policy",
		IsAuthenticated: isAuthenticated(c),
		IsAdmin:         isAdmin(c),
		CacheParameter:  controller.config.CacheParameter,
	}
	c.HTML(http.StatusOK, "privacy-policy.html", pd)
}
