// Package routes defines all the handling functions for all the routes
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/tournify/web/config"
	"github.com/tournify/web/lang"
	"github.com/tournify/web/middleware"
	"github.com/tournify/web/models"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Controller struct {
	db     *gorm.DB
	bundle *i18n.Bundle
	config config.Config
}

func New(db *gorm.DB, c config.Config, b *i18n.Bundle) Controller {
	return Controller{
		db:     db,
		config: c,
		bundle: b,
	}
}

type PageData struct {
	Trans             func(s string) string
	Title             string
	SiteName          string
	Messages          []Message
	IsAuthenticated   bool
	IsAdmin           bool
	CacheParameter    string
	CreateTournament  string
	Dashboard         string
	Blog              string
	CreateBlogPost    string
	Admin             string
	Logout            string
	Login             string
	Register          string
	Search            string
	Year              string
	ForkThisProjectOn string
	CreatedBy         string
	PrivacyPolicy     string
	TermsOfService    string
}

type Message struct {
	Type    string // success, warning, error, etc.
	Content string
}

func isAuthenticated(c *gin.Context) bool {
	_, exists := c.Get(middleware.UserIDKey)
	return exists
}

func isUnauthenticatedSession(c *gin.Context) bool {
	_, userIDExists := c.Get(middleware.UserIDKey)
	_, sessionIDExists := c.Get(middleware.SessionIDKey)
	return !userIDExists && sessionIDExists
}

func isAdmin(c *gin.Context) bool {
	_, exists := c.Get(middleware.UserIDKey)
	role, roleExists := c.Get(middleware.UserRoleKey)
	return exists && roleExists && role == "admin"
}

func canEditTournament(c *gin.Context, tournamentID uint) bool {
	userTournaments, userTournamentsExists := c.Get(middleware.UserTournamentsKey)
	if userTournamentsExists && userTournaments != nil {
		if _, ok := userTournaments.([]models.Tournament); ok {
			for _, t := range userTournaments.([]models.Tournament) {
				if t.ID == tournamentID {
					return true
				}
			}
		}
	} else {
		sessionTournaments, sessionTournamentsExists := c.Get(middleware.SessionTournamentsKey)
		if sessionTournamentsExists && sessionTournaments != nil {
			if _, ok := sessionTournaments.([]models.Tournament); ok {
				for _, t := range sessionTournaments.([]models.Tournament) {
					if t.ID == tournamentID {
						return true
					}
				}
			}
		}
	}
	return false
}

func domainLanguage(c *gin.Context) string {
	if strings.Contains(c.Request.Host, "turnering.io") {
		return "se"
	}
	return "en"
}

func (controller Controller) defaultPageData(c *gin.Context) PageData {
	langService := lang.New(c, controller.bundle)

	return PageData{
		Title:           langService.Trans("Home"),
		Year:            time.Now().Format("2006"),
		IsAuthenticated: isAuthenticated(c),
		IsAdmin:         isAdmin(c),
		CacheParameter:  controller.config.CacheParameter,
		Trans:           langService.Trans,
	}
}
