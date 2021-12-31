package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/tournify/web/models"
	"log"
	"net/http"
)

type SearchData struct {
	PageData
	Results []models.Tournament
}

func (controller Controller) Search(c *gin.Context) {
	localize := i18n.NewLocalizer(controller.bundle, domainLanguage(c))

	title, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "search_title",
			Other: "Search",
		},
	})
	pd := controller.defaultPageData(c)
	pd.Title = title
	sd := SearchData{
		PageData: pd,
	}
	search := c.PostForm("search")

	var results []models.Tournament

	search = fmt.Sprintf("%s%s%s", "%", search, "%")

	res := controller.db.Where("privacy = ? AND (name LIKE ? OR description LIKE ?)", models.TournamentPrivacyPublic, search, search).Limit(100).Find(&results)

	if res.Error != nil || len(results) == 0 {
		pd.Messages = append(sd.Messages, Message{
			Type:    "error",
			Content: "No results found",
		})
		log.Println(res.Error)
		c.HTML(http.StatusOK, "search.html", sd)
		return
	}

	sd.Results = results

	c.HTML(http.StatusOK, "search.html", sd)
}
