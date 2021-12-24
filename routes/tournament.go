package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type TournamentCreatePageData struct {
	PageData
	TournamentName string
}

func (controller Controller) TournamentCreate(c *gin.Context) {
	pd := TournamentCreatePageData{
		PageData: PageData{
			Title:           "Create Tournament",
			IsAuthenticated: isAuthenticated(c),
			CacheParameter:  controller.config.CacheParameter,
		},
		TournamentName: "",
	}
	c.HTML(http.StatusOK, "tournament-create.html", pd)
}

func (controller Controller) TournamentCreatePost(c *gin.Context) {
	pd := TournamentCreatePageData{
		PageData: PageData{
			Title:           "Create Tournament",
			IsAuthenticated: isAuthenticated(c),
			CacheParameter:  controller.config.CacheParameter,
		},
		TournamentName: "",
	}
	pd.TournamentName = c.PostForm("tourname")
	tourType := c.PostForm("tourtype")
	if tourType == "" {
		// Submit from index page, do not generate a tournament
		c.HTML(http.StatusOK, "tournament-create.html", pd)
		return
	}

	// TODO generate a tournament

}

func (controller Controller) TournamentView(context *gin.Context) {

}

func (controller Controller) TournamentViewPost(context *gin.Context) {

}
