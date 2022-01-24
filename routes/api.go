package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tournify/tournify"
	"github.com/tournify/web/lang"
	"github.com/tournify/web/models"
	"gorm.io/gorm"
	"log"
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
		log.Println(err)
		c.JSON(http.StatusBadRequest, APIError{
			Error: "Could not read the incoming data.",
		})
		return
	}

	idInt, err := strconv.Atoi(req.ID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, APIError{
			Error: "Could not read the incoming data.",
		})
		return
	}

	homeInt, err2 := strconv.Atoi(req.Home)
	if err2 != nil {
		log.Println(err2)
		c.JSON(http.StatusBadRequest, APIError{
			Error: "Could not read the incoming data.",
		})
		return
	}

	awayInt, err3 := strconv.Atoi(req.Away)
	if err3 != nil {
		log.Println(err3)
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

	res := controller.db.Where(t).Preload("Games.Parents.Teams").First(&t)
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

	if len(g.Teams) < 2 {
		c.JSON(http.StatusNotFound, APIError{
			Error: "Can not set score for requested game.",
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

	// For elimination tournaments we want to update the games in the next round
	if t.Type == 1 {
		for _, tg := range t.Games {
			for _, parent := range tg.Parents {
				if parent.GetID() == g.GetID() {
					var winTeam models.Team
					var loseTeam models.Team
					winTeamTmp := GetWinnerTeam(g)
					loseTeamTmp := GetLoserTeam(g)
					// Win team is nil in a tie game
					if winTeamTmp != nil {
						if len(tg.Teams) == 0 {
							err = controller.db.Model(&tg).Association("Teams").Find(&tg.Teams)
							if err != nil {
								log.Println(err)
								c.JSON(http.StatusInternalServerError, APIError{
									Error: "Could not save score",
								})
								return
							}
						}
						for ti := range parent.Teams {
							if parent.Teams[ti].GetID() == winTeamTmp.GetID() {
								winTeam = parent.Teams[ti]
							} else if parent.Teams[ti].GetID() == loseTeamTmp.GetID() {
								loseTeam = parent.Teams[ti]
							}
						}

						if tg.HomeTeamID == nil && (tg.AwayTeamID == nil || loseTeam.GetID() != int(*tg.AwayTeamID)) {
							tg.SetHomeTeam(&winTeam)
						} else if tg.AwayTeamID == nil && (tg.HomeTeamID == nil || loseTeam.GetID() != int(*tg.HomeTeamID)) {
							tg.SetAwayTeam(&winTeam)
						} else if tg.HomeTeamID != nil && loseTeam.GetID() == int(*tg.HomeTeamID) && (tg.AwayTeamID == nil || loseTeam.GetID() != int(*tg.AwayTeamID)) {
							tg.SetHomeTeam(&winTeam)
						} else if tg.AwayTeamID != nil && loseTeam.GetID() == int(*tg.AwayTeamID) && (tg.HomeTeamID == nil || loseTeam.GetID() != int(*tg.HomeTeamID)) {
							tg.SetAwayTeam(&winTeam)
						}
						loseTeam.SetEliminatedCount(1)
						winTeam.SetEliminatedCount(0)

						res = controller.db.Save(winTeam)
						if res.Error != nil {
							c.JSON(http.StatusInternalServerError, APIError{
								Error: "Could not save score",
							})
							return
						}

						res = controller.db.Save(loseTeam)
						if res.Error != nil {
							c.JSON(http.StatusInternalServerError, APIError{
								Error: "Could not save score",
							})
							return
						}

						// Fix for teams
						var teams []models.Team
						for i := range tg.Teams {
							if tg.Teams[i].GetID() != -1 {
								teams = append(teams, tg.Teams[i])
							}
						}
						tg.Teams = teams
						res = controller.db.Save(tg)
						if res.Error != nil {
							c.JSON(http.StatusInternalServerError, APIError{
								Error: "Could not save score",
							})
							return
						}
						err = controller.db.Model(&tg).Association("Teams").Replace(tg.Teams)
						if err != nil {
							log.Println(err)
							c.JSON(http.StatusInternalServerError, APIError{
								Error: "Could not save score",
							})
							return
						}
					}
				}
			}
		}
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

func (controller Controller) APITournamentGames(c *gin.Context) {
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

	if t.Type == 0 {
		var games []models.Game
		game := models.Game{
			TournamentID: t.ID,
		}

		res = controller.db.Where(game).Preload("Teams").Preload("Scores").Order("depth DESC").Find(&games)
		if res.Error != nil {
			c.JSON(http.StatusNotFound, APIError{
				Error: "Could not find games",
			})
			return
		}
		c.JSON(http.StatusOK, games)
		return
	} else if t.Type == 1 {
		langService := lang.New(c, controller.bundle)
		rounds, err := controller.getEliminationTournamentGames(t, langService.Trans)
		if err != nil {
			c.JSON(http.StatusNotFound, APIError{
				Error: "Could not find games",
			})
			return
		}
		c.JSON(http.StatusOK, rounds)
		return
	}
	c.JSON(http.StatusNotFound, APIError{
		Error: "Could not find games",
	})
	return
}

func GetWinnerTeam(g models.Game) tournify.TeamInterface {
	if g.GetHomeTeam().GetID() == -1 {
		return g.GetAwayTeam()
	} else if g.GetAwayTeam().GetID() == -1 {
		return g.GetHomeTeam()
	} else if g.GetAwayTeam().GetID() == -1 && g.GetHomeTeam().GetID() == -1 {
		return nil
	}
	if g.GetAwayScore().GetPoints() > g.GetHomeScore().GetPoints() {
		return g.GetAwayTeam()
	} else if g.GetHomeScore().GetPoints() > g.GetAwayScore().GetPoints() {
		return g.GetHomeTeam()
	}
	return nil
}

func GetLoserTeam(g models.Game) tournify.TeamInterface {
	if winTeam := GetWinnerTeam(g); winTeam != nil {
		if winTeam.GetID() == g.GetAwayTeam().GetID() {
			return g.GetHomeTeam()
		} else {
			return g.GetAwayTeam()
		}
	}
	return nil
}
