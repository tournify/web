package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/tournify/tournify"
	"github.com/tournify/web/models"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TournamentCreatePageData struct {
	PageData
	TournamentName string
}

type TournamentViewPageData struct {
	PageData
	TournamentName string
	TournamentSlug string
	Groups         map[int]TournamentViewGroup
}

type TournamentViewGroup struct {
	Stats []models.Statistics `json:"stats"`
	Games []models.Game       `json:"-"`
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
	name := strings.TrimSpace(c.PostForm("tourname"))

	pd.TournamentName = name
	tourType := c.PostForm("tourtype")
	if tourType == "" {
		// Submit from index page, do not generate a tournament
		c.HTML(http.StatusOK, "tournament-create.html", pd)
		return
	}

	if tourType == "0" {
		meetCount := c.PostForm("meetcount")
		groupCount := c.PostForm("groupcount")
		elimCount := c.PostForm("elimcount")
		winPoints := c.PostForm("winpoints")
		tiePoints := c.PostForm("tiepoints")
		lossPoints := c.PostForm("losspoints")
		visibility := c.PostForm("visibility")
		teams := c.PostFormArray("team[]")
		meetCountInt, err := strconv.Atoi(meetCount)
		if err != nil {
			log.Println(err)
			pd.Messages = append(pd.Messages, Message{
				Type:    "error",
				Content: "Bad input, please check your settings and try again.",
			})
			c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
			return
		}
		groupCountInt, err2 := strconv.Atoi(groupCount)
		if err2 != nil {
			log.Println(err2)
			pd.Messages = append(pd.Messages, Message{
				Type:    "error",
				Content: "Bad input, please check your settings and try again.",
			})
			c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
			return
		}
		_, err3 := strconv.Atoi(elimCount)
		if err3 != nil {
			log.Println(err3)
			pd.Messages = append(pd.Messages, Message{
				Type:    "error",
				Content: "Bad input, please check your settings and try again.",
			})
			c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
			return
		}
		_, err4 := strconv.Atoi(winPoints)
		if err4 != nil {
			log.Println(err4)
			pd.Messages = append(pd.Messages, Message{
				Type:    "error",
				Content: "Bad input, please check your settings and try again.",
			})
			c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
			return
		}
		_, err5 := strconv.Atoi(tiePoints)
		if err5 != nil {
			log.Println(err5)
			pd.Messages = append(pd.Messages, Message{
				Type:    "error",
				Content: "Bad input, please check your settings and try again.",
			})
			c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
			return
		}
		_, err6 := strconv.Atoi(lossPoints)
		if err6 != nil {
			log.Println(err6)
			pd.Messages = append(pd.Messages, Message{
				Type:    "error",
				Content: "Bad input, please check your settings and try again.",
			})
			c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
			return
		}
		visibilityInt, err7 := strconv.Atoi(visibility)
		if err7 != nil {
			log.Println(err7)
			pd.Messages = append(pd.Messages, Message{
				Type:    "error",
				Content: "Bad input, please check your settings and try again.",
			})
			c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
			return
		}

		// TODO change this with paid users
		if len(teams) > 10 {
			pd.Messages = append(pd.Messages, Message{
				Type:    "error",
				Content: "Up to 10 teams is currently supported.",
			})
			c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
			return
		}

		if len(teams) < 2 {
			pd.Messages = append(pd.Messages, Message{
				Type:    "error",
				Content: "At least 2 teams have to be in a tournament",
			})
			c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
			return
		}

		if meetCountInt > 4 {
			pd.Messages = append(pd.Messages, Message{
				Type:    "error",
				Content: "Up to a meeting count of 4 is currently supported.",
			})
			c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
			return
		}

		if meetCountInt < 1 {
			pd.Messages = append(pd.Messages, Message{
				Type:    "error",
				Content: "Meet count has to be at least 1.",
			})
			c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
			return
		}

		if groupCountInt < 1 {
			pd.Messages = append(pd.Messages, Message{
				Type:    "error",
				Content: "Group count has to be at least 1.",
			})
			c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
			return
		}

		if name == "" {
			name = "Tournament " + time.Now().Format("2006-01-02 15:04:05")
		}
		slugString := slug.Make(name)
		tournamentModel := models.Tournament{
			Name:    name,
			Slug:    controller.createUniqueTournamentSlug(slugString, 0),
			Type:    0,
			Privacy: visibilityInt,
		}
		res := controller.db.Save(&tournamentModel)
		if res.Error != nil {
			pd.Messages = append(pd.Messages, Message{
				Type:    "error",
				Content: "Could not create tournament, please try again or contact support.",
			})
			c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
			return
		}

		eliminationCountOption := models.TournamentOption{
			Key:          "elimination_count",
			Value:        elimCount,
			TournamentID: tournamentModel.ID,
		}
		res = controller.db.Save(&eliminationCountOption)
		if res.Error != nil {
			pd.Messages = append(pd.Messages, Message{
				Type:    "error",
				Content: "Could not create tournament, please try again or contact support.",
			})
			c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
			return
		}
		winPointsOption := models.TournamentOption{
			Key:          "win_points",
			Value:        winPoints,
			TournamentID: tournamentModel.ID,
		}
		res = controller.db.Save(&winPointsOption)
		if res.Error != nil {
			pd.Messages = append(pd.Messages, Message{
				Type:    "error",
				Content: "Could not create tournament, please try again or contact support.",
			})
			c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
			return
		}
		lossPointsOption := models.TournamentOption{
			Key:          "loss_points",
			Value:        lossPoints,
			TournamentID: tournamentModel.ID,
		}
		res = controller.db.Save(&lossPointsOption)
		if res.Error != nil {
			pd.Messages = append(pd.Messages, Message{
				Type:    "error",
				Content: "Could not create tournament, please try again or contact support.",
			})
			c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
			return
		}
		tiePointsOption := models.TournamentOption{
			Key:          "tie_points",
			Value:        tiePoints,
			TournamentID: tournamentModel.ID,
		}
		res = controller.db.Save(&tiePointsOption)
		if res.Error != nil {
			pd.Messages = append(pd.Messages, Message{
				Type:    "error",
				Content: "Could not create tournament, please try again or contact support.",
			})
			c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
			return
		}

		var teamModels []models.Team
		teamCount := 0
		for _, team := range teams {
			teamCount++
			if team == "" {
				team = fmt.Sprintf("Team %d", teamCount)
			}
			teamSlug := slug.Make(team)
			teamModel := models.Team{
				Name: team,
				Slug: controller.createUniqueTeamSlug(teamSlug, 0),
			}
			teamModels = append(teamModels, teamModel)
		}

		res = controller.db.CreateInBatches(teamModels, 100)
		if res.Error != nil {
			pd.Messages = append(pd.Messages, Message{
				Type:    "error",
				Content: "Could not create teams, please try again or contact support.",
			})
			c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
			return
		}

		teamInterfaces := make([]tournify.TeamInterface, len(teams))

		for i := range teams {
			teamInterfaces[i] = &teamModels[i]
		}

		// The CreateGroupTournamentFromTeams method takes a slice of teams along with the group count and meet count
		tournament := tournify.CreateGroupTournamentFromTeams(teamInterfaces, groupCountInt, meetCountInt)

		for i, group := range tournament.GetGroups() {
			groupName := fmt.Sprintf("Group %d", i+1)
			groupSlug := slug.Make(groupName)
			groupModel := models.Group{
				Name:         groupName,
				Slug:         controller.createUniqueGroupSlug(groupSlug, 0),
				TournamentID: tournamentModel.ID,
			}
			res = controller.db.Save(&groupModel)
			if res.Error != nil {
				pd.Messages = append(pd.Messages, Message{
					Type:    "error",
					Content: "Could not create groups, please try again or contact support.",
				})
				c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
				return
			}
			if group.GetGames() != nil {
				for x, game := range *group.GetGames() {
					var gameTeams []models.Team
					for _, t := range game.GetTeams() {
						gameTeams = append(gameTeams, models.Team{
							Model: gorm.Model{
								ID: uint(t.GetID()),
							},
						})
					}
					gameName := fmt.Sprintf("Game %d", x+1)
					gameSlug := slug.Make(gameName)
					gameModel := models.Game{
						Name:         gameName,
						Slug:         controller.createUniqueGameSlug(gameSlug, 0),
						TournamentID: tournamentModel.ID,
						GroupID:      &groupModel.ID,
						Teams:        gameTeams,
					}
					res = controller.db.Save(&gameModel)
					if res.Error != nil {
						pd.Messages = append(pd.Messages, Message{
							Type:    "error",
							Content: "Could not create games, please try again or contact support.",
						})
						c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
						return
					}
				}
			}
		}
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/tournament/%s", tournamentModel.Slug))
		return
	}
	pd.Messages = append(pd.Messages, Message{
		Type:    "error",
		Content: "Tournament type not implemented.",
	})
	c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
}

func (controller Controller) TournamentView(c *gin.Context) {
	pd := TournamentViewPageData{
		PageData: PageData{
			Title:           "Tournament",
			IsAuthenticated: isAuthenticated(c),
			CacheParameter:  controller.config.CacheParameter,
		},
		TournamentName: "",
	}
	slugParam := c.Param("slug")
	t := models.Tournament{
		Slug: slugParam,
	}
	// TODO handle privacy here
	res := controller.db.Where(t).First(&t)
	if res.Error != nil {
		c.HTML(http.StatusNotFound, "404.html", pd)
		return
	}

	pd.TournamentSlug = t.Slug
	pd.TournamentName = t.Name

	var err error
	pd.Groups, err = controller.getGroupTournamentStats(t)
	if err != nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: "Could not generate statistics, please try again or contact support.",
		})
		c.HTML(http.StatusBadRequest, "tournament.html", pd)
		return
	}

	c.HTML(http.StatusOK, "tournament.html", pd)
}

func hasTeam(group tournify.Group, team models.Team) bool {
	for _, t := range group.Teams {
		if t.GetID() == team.GetID() {
			return true
		}
	}
	return false
}

func (controller *Controller) getGroupTournamentStats(t models.Tournament) (map[int]TournamentViewGroup, error) {
	var games []models.Game
	game := models.Game{
		TournamentID: t.ID,
	}

	res := controller.db.Where(game).Preload("Teams").Preload("Scores").Find(&games)
	if res.Error != nil {
		return nil, res.Error
	}

	var options []models.TournamentOption
	option := models.TournamentOption{
		TournamentID: t.ID,
	}

	res = controller.db.Where(option).Find(&options)
	if res.Error != nil {
		return nil, res.Error
	}

	winPoints := 3
	lossPoints := 0
	tiePoints := 1

	for _, o := range options {
		vInt, err8 := strconv.Atoi(o.Value)
		if err8 == nil {
			if o.Key == "win_points" {
				winPoints = vInt
			} else if o.Key == "loss_points" {
				lossPoints = vInt
			} else if o.Key == "tie_points" {
				tiePoints = vInt
			}
		}
	}

	var groups []tournify.Group
	var groupInterfaces []tournify.GroupInterface
	viewGroups := map[int]TournamentViewGroup{}

	for gi, g := range games {
		if _, ok := viewGroups[int(*g.GroupID)]; !ok {
			viewGroups[int(*g.GroupID)] = TournamentViewGroup{}
		}
		tmpViewGroup := viewGroups[int(*g.GroupID)]
		tmpViewGroup.Games = append(tmpViewGroup.Games, g)
		viewGroups[int(*g.GroupID)] = tmpViewGroup

		groupIndex := -1
		for i, group := range groups {
			if group.GetID() == int(*g.GroupID) {
				groupIndex = i
			}
		}
		if groupIndex == -1 {
			group := tournify.Group{
				ID: int(*g.GroupID),
			}
			if !g.CreatedAt.Equal(g.UpdatedAt) {
				group.Games = append(group.Games, &games[gi])
			}
			for i, team := range g.Teams {
				if !hasTeam(group, team) {
					group.Teams = append(group.Teams, &games[gi].Teams[i])
				}
			}
			groups = append(groups, group)
		} else {
			if !g.CreatedAt.Equal(g.UpdatedAt) {
				groups[groupIndex].Games = append(groups[groupIndex].Games, &games[gi])
			}
			for i, team := range g.Teams {
				if !hasTeam(groups[groupIndex], team) {
					groups[groupIndex].Teams = append(groups[groupIndex].Teams, &games[gi].Teams[i])
				}
			}
		}
	}

	for i := range groups {
		groupInterfaces = append(groupInterfaces, &groups[i])
	}

	tournament := tournify.Tournament{
		Type:   tournify.TournamentTypeGroup,
		Groups: groupInterfaces,
	}

	stats, err := tournify.GetGroupTournamentStats(&tournament, winPoints, lossPoints, tiePoints)
	if err != nil {
		return nil, err
	}

	for _, stat := range stats {
		if _, ok := viewGroups[stat.GetGroup().GetID()]; !ok {
			viewGroups[stat.GetGroup().GetID()] = TournamentViewGroup{}
		}
		tmpViewGroup := viewGroups[stat.GetGroup().GetID()]
		statistic := models.Statistics{
			Group:         stat.GetGroup(),
			Team:          *stat.GetTeam().(*models.Team),
			Played:        stat.GetPlayed(),
			Wins:          stat.GetWins(),
			Losses:        stat.GetLosses(),
			Ties:          stat.GetTies(),
			PointsFor:     stat.GetPointsFor(),
			PointsAgainst: stat.GetPointsAgainst(),
			Points:        stat.GetPoints(),
		}
		tmpViewGroup.Stats = append(tmpViewGroup.Stats, statistic)
		viewGroups[stat.GetGroup().GetID()] = tmpViewGroup
	}

	// Normalize group ids

	normalizedViewGroups := map[int]TournamentViewGroup{}
	count := 0
	for _, vg := range viewGroups {
		count++
		normalizedViewGroups[count] = vg
	}

	return normalizedViewGroups, nil
}
