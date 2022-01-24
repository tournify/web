package web

import (
	"embed"
	"github.com/BurntSushi/toml"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/tournify/web/middleware"
	"github.com/tournify/web/routes"
	"golang.org/x/text/language"
	"html/template"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"time"
)

//go:embed dist/*
var staticFS embed.FS

func Run() {
	// When generating random strings we need to provide a seed otherwise we always get the same strings the next time our application starts
	rand.Seed(time.Now().UnixNano())

	// Translations
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	_, err := bundle.LoadMessageFile("active.en.toml")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = bundle.LoadMessageFile("active.se.toml")
	if err != nil {
		log.Fatalln(err)
	}

	var t *template.Template
	conf := loadEnvVariables()

	db, err2 := connectToDatabase(conf)
	if err2 != nil {
		log.Fatalln(err2)
	}

	err = migrateDatabase(db)
	if err != nil {
		log.Fatalln(err)
	}

	t, err = loadTemplates()
	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()

	store := cookie.NewStore([]byte(conf.CookieSecret))
	r.Use(sessions.Sessions("tournify_session", store))

	r.SetHTMLTemplate(t)

	subFS, err := fs.Sub(staticFS, "dist/assets")
	if err != nil {
		log.Fatalln(err)
	}

	assets := r.Group("/assets")
	assets.Use(middleware.Cache())

	assets.StaticFS("/", http.FS(subFS))

	r.Use(middleware.Session(db))
	r.Use(middleware.General())

	controller := routes.New(db, conf, bundle)

	r.GET("/", controller.Index)
	r.GET("/search", controller.Search)
	r.POST("/search", controller.Search)
	r.GET("/terms-of-service", controller.TermsOfService)
	r.GET("/privacy-policy", controller.PrivacyPolicy)
	r.GET("/tournament/create", controller.TournamentCreate)
	r.POST("/tournament/create", controller.TournamentCreatePost)
	r.Any("/tournament/:slug", controller.TournamentView)
	r.GET("/tournament/:slug/game/:gameslug", controller.TournamentGameView)
	r.GET("/blog", controller.BlogView)
	r.GET("/blog/:slug", controller.BlogViewPage)
	r.POST("/subscribe", controller.SubscribePost)

	api := r.Group("/api")
	api.Use(middleware.Sensitive())
	api.POST("/tournament/:slug/game/:id", controller.APITournamentGameUpdate)
	api.GET("/tournament/:slug/stats", controller.APITournamentStats)
	api.GET("/tournament/:slug/games", controller.APITournamentGames)

	r.NoRoute(controller.NoRoute)

	noAuth := r.Group("/")
	noAuth.Use(middleware.NoAuth())

	noAuth.GET("/login", controller.Login)
	noAuth.GET("/register", controller.Register)
	noAuth.GET("/activate/resend", controller.ResendActivation)
	noAuth.GET("/activate/:token", controller.Activate)
	noAuth.GET("/user/password/forgot", controller.ForgotPassword)
	noAuth.GET("/user/password/reset/:token", controller.ResetPassword)

	noAuthPost := noAuth.Group("/")
	noAuthPost.Use(middleware.Throttle(conf.RequestsPerMinute))

	noAuthPost.POST("/login", controller.LoginPost)
	noAuthPost.POST("/register", controller.RegisterPost)
	noAuthPost.POST("/activate/resend", controller.ResendActivationPost)
	noAuthPost.POST("/user/password/forgot", controller.ForgotPasswordPost)
	noAuthPost.POST("/user/password/reset/:token", controller.ResetPasswordPost)

	auth := r.Group("/")
	auth.Use(middleware.Auth())
	auth.Use(middleware.Sensitive())

	auth.GET("/auth", controller.Admin)
	auth.POST("/auth", controller.Admin)
	auth.GET("/logout", controller.Logout)
	auth.GET("/dashboard", controller.NoRoute)

	admin := auth.Group("/")
	admin.Use(middleware.Admin())
	admin.GET("/admin", controller.NoRoute)
	admin.GET("/blog/create", controller.BlogCreate)
	admin.POST("/blog/create", controller.BlogCreatePost)

	err = r.Run(conf.Port)
	if err != nil {
		log.Fatalln(err)
	}
}
