package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tournify/web/models"
	"log"
	"net/http"
)

type IndexData struct {
	PageData
	RecentTournaments []models.Tournament
}

func (controller Controller) Index(c *gin.Context) {

	// Select 10 most recent tournaments
	var recentTournaments []models.Tournament
	res := controller.db.Where("privacy = ?", models.TournamentPrivacyPublic).Limit(10).Order("created_at DESC").Find(&recentTournaments)
	if res.Error != nil {
		log.Println(res.Error)
		// We just load the index page instead of showing an error if this fails since the index page is critical but this section is not
	}

	id := IndexData{
		PageData:          controller.defaultPageData(c),
		RecentTournaments: recentTournaments,
	}
	c.HTML(http.StatusOK, "index.html", id)
}
