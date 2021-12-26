// Package routes defines all the handling functions for all the routes
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/tournify/web/config"
	"github.com/tournify/web/middleware"
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

func isAdmin(c *gin.Context) bool {
	_, exists := c.Get(middleware.UserIDKey)
	role, roleExists := c.Get(middleware.UserRoleKey)
	return exists && roleExists && role == "admin"
}

func domainLanguage(c *gin.Context) string {
	if strings.Contains(c.Request.Host, "turnering.io") {
		return "se"
	}
	return "en"
}

func (controller Controller) defaultPageData(c *gin.Context) PageData {
	localize := i18n.NewLocalizer(controller.bundle, domainLanguage(c))

	home, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "home",
			Other: "Home",
		},
	})

	siteName, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "site_name",
			Other: "Tournify.io",
		},
	})

	createTournament, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "create_tournament",
			Other: "Create Tournament",
		},
	})

	dashboard, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "dashboard",
			Other: "Dashboard",
		},
	})

	blog, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "blog",
			Other: "Blog",
		},
	})

	createBlogPost, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "create_blog_post",
			Other: "Create Blog Post",
		},
	})

	admin, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "admin",
			Other: "Admin",
		},
	})

	logout, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "logout",
			Other: "Logout",
		},
	})

	login, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "login",
			Other: "Login",
		},
	})

	register, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "register",
			Other: "Register",
		},
	})

	search, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "search",
			Other: "Search",
		},
	})

	forkThisProjectOn, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "fork_this_project_on",
			Other: "Fork this project on",
		},
	})

	createdBy, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "created_by",
			Other: "Created by",
		},
	})

	privacyPolicy, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "privacy_policy",
			Other: "Privacy Policy",
		},
	})

	termsOfService, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "terms_of_service",
			Other: "Terms of Service",
		},
	})

	return PageData{
		Title:             home,
		SiteName:          siteName,
		CreateTournament:  createTournament,
		Dashboard:         dashboard,
		Blog:              blog,
		CreateBlogPost:    createBlogPost,
		Admin:             admin,
		Logout:            logout,
		Login:             login,
		Register:          register,
		Search:            search,
		ForkThisProjectOn: forkThisProjectOn,
		CreatedBy:         createdBy,
		PrivacyPolicy:     privacyPolicy,
		TermsOfService:    termsOfService,
		Year:              time.Now().Format("2006"),
		IsAuthenticated:   isAuthenticated(c),
		IsAdmin:           isAdmin(c),
		CacheParameter:    controller.config.CacheParameter,
	}
}
