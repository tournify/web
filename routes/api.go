package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tournify/web/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type APIResponse struct {
	Message string `json:"message"`
}

type APIError struct {
	Error string `json:"error"`
}

type APITournamentGameUpdateRequest struct {
	Away string `json:"away"`
	Home string `json:"home"`
	ID   string `json:"id"`
	Slug string `json:"slug"`
}

func (controller Controller) APITournamentGameUpdate(c *gin.Context) {
	slugParam := c.Param("slug")
	gameIdParam := c.Param("id")
	// TODO check session permission here

	var req APITournamentGameUpdateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIError{
			Error: "Could not read the incoming data.",
		})
		return
	}

	idInt, err := strconv.Atoi(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIError{
			Error: "Could not read the incoming data.",
		})
		return
	}

	homeInt, err2 := strconv.Atoi(req.Home)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, APIError{
			Error: "Could not read the incoming data.",
		})
		return
	}

	awayInt, err3 := strconv.Atoi(req.Away)
	if err3 != nil {
		c.JSON(http.StatusBadRequest, APIError{
			Error: "Could not read the incoming data.",
		})
		return
	}

	if slugParam != req.Slug || gameIdParam != req.ID {
		c.JSON(http.StatusBadRequest, APIError{
			Error: "Data is invalid.",
		})
		return
	}

	t := models.Tournament{
		Slug: slugParam,
	}

	res := controller.db.Where(t).First(&t)
	if res.Error != nil {
		c.JSON(http.StatusNotFound, APIError{
			Error: "Could not find the requested tournament.",
		})
		return
	}

	if !canEditTournament(c, t.ID) {
		c.JSON(http.StatusUnauthorized, APIError{
			Error: "Could not update the requested tournament.",
		})
		return
	}

	g := models.Game{
		TournamentID: t.ID,
	}
	g.ID = uint(idInt)

	res = controller.db.Where(g).Preload("Scores").Preload("Teams").First(&g)
	if res.Error != nil {
		c.JSON(http.StatusNotFound, APIError{
			Error: "Could not find the requested game.",
		})
		return
	}

	g.SetScore(float64(homeInt), float64(awayInt))

	res = controller.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(g)

	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, APIError{
			Error: "Could not save score",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Message: "Game saved successfully",
	})
}

func (controller Controller) APITournamentStats(c *gin.Context) {
	slugParam := c.Param("slug")
	t := models.Tournament{
		Slug: slugParam,
	}
	// TODO handle privacy here
	res := controller.db.Where(t).First(&t)
	if res.Error != nil {
		c.JSON(http.StatusNotFound, APIError{
			Error: "Could not find tournament",
		})
		return
	}

	groups, err := controller.getGroupTournamentStats(t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIError{
			Error: "Could not generate stats",
		})
		return
	}
	c.JSON(http.StatusOK, groups)
}
