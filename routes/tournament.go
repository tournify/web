package routes

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/tournify/tournify"
	"github.com/tournify/web/middleware"
	"github.com/tournify/web/models"
	"github.com/tournify/web/util"
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
	TournamentType int
	Groups         map[int]TournamentViewGroup
	Rounds         []Round `json:"rounds"`
}

type Round struct {
	Games []models.Game `json:"games"`
	Name  string        `json:"name"`
}

type TournamentViewGroup struct {
	Stats []models.Statistics `json:"stats"`
	Games []models.Game       `json:"-"`
}

func (controller Controller) TournamentCreate(c *gin.Context) {
	localize := i18n.NewLocalizer(controller.bundle, domainLanguage(c))

	title, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "create_tournament_title",
			Other: "Create Tournament",
		},
	})
	pd := TournamentCreatePageData{
		PageData:       controller.defaultPageData(c),
		TournamentName: "",
	}
	pd.Title = title
	c.HTML(http.StatusOK, "tournament-create.html", pd)
}

func (controller Controller) TournamentCreatePost(c *gin.Context) {
	localize := i18n.NewLocalizer(controller.bundle, domainLanguage(c))

	title, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "create_tournament_title",
			Other: "Create Tournament",
		},
	})
	pd := TournamentCreatePageData{
		PageData:       controller.defaultPageData(c),
		TournamentName: "",
	}
	pd.Title = title
	name := strings.TrimSpace(c.PostForm("tourname"))

	pd.TournamentName = name
	tourType := c.PostForm("tourtype")
	if tourType == "" {
		// Submit from index page, do not generate a tournament
		c.HTML(http.StatusOK, "tournament-create.html", pd)
		return
	}

	visibility := c.PostForm("visibility")
	teams := c.PostFormArray("team[]")

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

	tourTypeInt, err8 := strconv.Atoi(tourType)
	if err8 != nil {
		log.Println(err8)
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: "Bad input, please check your settings and try again.",
		})
		c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
		return
	}

	if len(teams) > 64 {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: "Up to 64 teams is currently supported.",
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

	if tourType == "0" {
		// Create group tournament
		meetCount := c.PostForm("meetcount")
		groupCount := c.PostForm("groupcount")
		winPoints := c.PostForm("winpoints")
		tiePoints := c.PostForm("tiepoints")
		lossPoints := c.PostForm("losspoints")
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

		if len(teams)/groupCountInt < 2 {
			pd.Messages = append(pd.Messages, Message{
				Type:    "error",
				Content: "There must be enough teams to have at least 2 teams per group.",
			})
			c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
			return
		}

		slugString, createError := controller.createTournament(c, tourTypeInt, name, visibilityInt, winPoints, lossPoints, tiePoints, groupCountInt, meetCountInt, teams, pd.Trans)
		if createError != nil {
			log.Println(createError)
			pd.Messages = append(pd.Messages, Message{
				Type:    "error",
				Content: "Could not create tournament, please try again or contact support.",
			})
			c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/tournament/%s", slugString))
		return
	} else if tourType == "1" {
		// Create elimination tournament
		slugString, createError := controller.createTournament(c, tourTypeInt, name, visibilityInt, "", "", "", 0, 0, teams, pd.Trans)
		if createError != nil {
			log.Println(createError)
			pd.Messages = append(pd.Messages, Message{
				Type:    "error",
				Content: "Could not create tournament, please try again or contact support.",
			})
			c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/tournament/%s", slugString))
		return
	}
	pd.Messages = append(pd.Messages, Message{
		Type:    "error",
		Content: "Tournament type not implemented.",
	})
	c.HTML(http.StatusBadRequest, "tournament-create.html", pd)
}

func (controller Controller) TournamentView(c *gin.Context) {
	localize := i18n.NewLocalizer(controller.bundle, domainLanguage(c))

	title, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "tournament_title",
			Other: "Tournament",
		},
	})
	pd := TournamentViewPageData{
		PageData:       controller.defaultPageData(c),
		TournamentName: "",
	}
	pd.Title = title
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
	pd.TournamentType = t.Type
	pd.CanEdit = canEditTournament(c, t.ID)

	if t.Type == 0 {
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
	} else if t.Type == 1 {
		var err error
		pd.Rounds, err = controller.getEliminationTournamentGames(t, pd.Trans)
		if err != nil {
			pd.Messages = append(pd.Messages, Message{
				Type:    "error",
				Content: "Could not generate statistics, please try again or contact support.",
			})
			c.HTML(http.StatusBadRequest, "tournament.html", pd)
			return
		}
	}

	c.HTML(http.StatusOK, "tournament.html", pd)
}

func (controller *Controller) getEliminationTournamentGames(t models.Tournament, trans func(s string) string) (rounds []Round, err error) {
	var games []models.Game
	game := models.Game{
		TournamentID: t.ID,
	}

	res := controller.db.Where(game).Preload("Teams").Preload("Scores").Order("depth DESC").Find(&games)
	if res.Error != nil {
		return nil, res.Error
	}
	lastDepth := -1
	index := 0
	for _, g := range games {
		if g.Depth != nil {
			if lastDepth == -1 {
				lastDepth = *g.Depth
				rounds = append(rounds, Round{
					Name: trans("Final"),
				})
			}
			if lastDepth != *g.Depth {
				lastDepth = *g.Depth
				index++
				if index == 1 {
					rounds = append(rounds, Round{
						Name: trans("Semi-Final"),
					})
				} else if index == 2 {
					rounds = append(rounds, Round{
						Name: trans("Quarter-Final"),
					})
				} else if index == 3 {
					rounds = append(rounds, Round{
						Name: trans("Eight-Final"),
					})
				} else if index == 4 {
					rounds = append(rounds, Round{
						Name: trans("16th-Final"),
					})
				} else if index == 5 {
					rounds = append(rounds, Round{
						Name: trans("32nd-Final"),
					})
				}
			}
			rounds[index].Games = append(rounds[index].Games, g)
		}
	}
	for i, j := 0, len(rounds)-1; i < j; i, j = i+1, j-1 {
		rounds[i], rounds[j] = rounds[j], rounds[i]
	}
	return rounds, nil
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
	normalizedViewGroups := normalizeViewGroups(viewGroups, map[int]TournamentViewGroup{}, 1)
	return normalizedViewGroups, nil
}

func normalizeViewGroups(gs map[int]TournamentViewGroup, normalized map[int]TournamentViewGroup, count int) map[int]TournamentViewGroup {
	if len(gs) == 0 {
		return normalized
	}
	lowestIndex := -1
	for index := range gs {
		if lowestIndex == -1 {
			lowestIndex = index
		} else if lowestIndex > index {
			lowestIndex = index
		}
	}
	normalized[count] = gs[lowestIndex]
	delete(gs, lowestIndex)
	return normalizeViewGroups(gs, normalized, count+1)
}

func (controller *Controller) createTournament(c *gin.Context, tourType int, name string, visibility int, winPoints string, lossPoints string, tiePoints string, groupCountInt int, meetCountInt int, teams []string, trans func(s string) string) (string, error) {
	if name == "" {
		name = "Tournament " + time.Now().Format("2006-01-02 15:04:05")
	}
	slugString := slug.Make(name)
	tournamentModel := models.Tournament{
		Name:    name,
		Slug:    controller.createUniqueTournamentSlug(slugString, 0),
		Type:    tourType,
		Privacy: visibility,
	}

	if isAuthenticated(c) {
		log.Println("authenticated user")
		userID, exists := c.Get(middleware.UserIDKey)

		if exists {
			userIDInt, ok := userID.(uint)
			if ok {
				userModel := models.User{}
				userModel.ID = userIDInt
				res := controller.db.Where(&userModel).First(&userModel)
				if res.Error != nil {
					log.Println(res.Error)
					return slugString, res.Error
				}

				// Associate the user to the tournament
				tournamentModel.Users = append(tournamentModel.Users, userModel)
			}
		} else {
			// This should never happen, but we log it so that we might see it if it does happen for some reason
			log.Println("userID doesn't exist but user is authenticated")
		}
	} else if isUnauthenticatedSession(c) {
		log.Println("unauthenticated session")
		sessionID, exists := c.Get(middleware.SessionIDKey)

		if exists {
			sessionIDInt, ok := sessionID.(uint)
			if ok {
				sessionModel := models.Session{}
				sessionModel.ID = sessionIDInt
				res := controller.db.Where(&sessionModel).First(&sessionModel)
				if res.Error != nil {
					log.Println(res.Error)
					return slugString, res.Error
				}

				// Associate the user to the tournament
				tournamentModel.Sessions = append(tournamentModel.Sessions, sessionModel)
			}
		} else {
			// This should never happen, but we log it so that we might see it if it does happen for some reason
			log.Println("userID doesn't exist but user is authenticated")
		}
	} else {
		log.Println("no session")
		// Generate a ULID for the current session
		sessionIdentifier := util.GenerateULID()

		ses := models.Session{
			Identifier: sessionIdentifier,
		}

		// Session is valid for 30 days
		ses.ExpiresAt = time.Now().Add(time.Hour * 24 * 30)

		res := controller.db.Save(&ses)
		if res.Error != nil {
			log.Println(res.Error)
			return slugString, res.Error
		}

		session := sessions.Default(c)
		session.Set(middleware.SessionIdentifierKey, sessionIdentifier)

		err := session.Save()
		if err != nil {
			log.Println(err)
			return slugString, err
		}

		tournamentModel.Sessions = append(tournamentModel.Sessions, ses)
	}

	res := controller.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&tournamentModel)
	if res.Error != nil {
		return slugString, res.Error
	}

	winPointsOption := models.TournamentOption{
		Key:          "win_points",
		Value:        winPoints,
		TournamentID: tournamentModel.ID,
	}
	res = controller.db.Save(&winPointsOption)
	if res.Error != nil {
		return slugString, res.Error
	}
	lossPointsOption := models.TournamentOption{
		Key:          "loss_points",
		Value:        lossPoints,
		TournamentID: tournamentModel.ID,
	}
	res = controller.db.Save(&lossPointsOption)
	if res.Error != nil {
		return slugString, res.Error
	}
	tiePointsOption := models.TournamentOption{
		Key:          "tie_points",
		Value:        tiePoints,
		TournamentID: tournamentModel.ID,
	}
	res = controller.db.Save(&tiePointsOption)
	if res.Error != nil {
		return slugString, res.Error
	}

	var teamModels []models.Team
	teamCount := 0
	for _, team := range teams {
		teamCount++
		teamSlug := slug.Make(team)
		if team == "" {
			team = fmt.Sprintf("%s %d", trans("Team"), teamCount)
			teamSlug = slug.Make(team) + "-" + util.RandomString(4)
		}
		teamModel := models.Team{
			Name: team,
			Slug: controller.createUniqueTeamSlug(teamSlug, 0),
		}
		teamModels = append(teamModels, teamModel)
	}

	res = controller.db.CreateInBatches(teamModels, 100)
	if res.Error != nil {
		return slugString, res.Error
	}

	teamInterfaces := make([]tournify.TeamInterface, len(teams))

	for i := range teams {
		teamInterfaces[i] = &teamModels[i]
	}

	if tourType == 0 {

		// The CreateGroupTournamentFromTeams method takes a slice of teams along with the group count and meet count
		tournament := tournify.CreateGroupTournamentFromTeams(teamInterfaces, groupCountInt, meetCountInt)

		for i, group := range tournament.GetGroups() {
			groupName := fmt.Sprintf("%s %d", trans("Group"), i+1)
			groupSlug := slug.Make(groupName) + "-" + util.RandomString(4)
			groupModel := models.Group{
				Name:         groupName,
				Slug:         controller.createUniqueGroupSlug(groupSlug, 0),
				TournamentID: tournamentModel.ID,
			}
			res = controller.db.Save(&groupModel)
			if res.Error != nil {
				return slugString, res.Error
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
					gameName := fmt.Sprintf("%s %d", trans("Game"), x+1)
					gameSlug := slug.Make(gameName) + "-" + util.RandomString(4)
					homeTeamID := uint(game.GetHomeTeam().GetID())
					awayTeamID := uint(game.GetAwayTeam().GetID())
					gameModel := models.Game{
						Name:         gameName,
						Slug:         controller.createUniqueGameSlug(gameSlug, 0),
						TournamentID: tournamentModel.ID,
						HomeTeamID:   &homeTeamID,
						AwayTeamID:   &awayTeamID,
						GroupID:      &groupModel.ID,
						Teams:        gameTeams,
					}
					res = controller.db.Save(&gameModel)
					if res.Error != nil {
						return slugString, res.Error
					}
				}
			}
		}
	} else if tourType == 1 {
		tournament := tournify.CreateEliminationTournamentFromTeams(teamInterfaces)
		if tournament.GetGames() != nil {
			insertedGames := make(map[int]int, len(tournament.GetGames()))
			for x, game := range tournament.GetGames() {
				var gameTeams []models.Team
				for _, t := range game.GetTeams() {
					if t.GetID() != -1 {
						gameTeams = append(gameTeams, models.Team{
							Model: gorm.Model{
								ID: uint(t.GetID()),
							},
						})
					}
				}
				depth := getGameDepth(game, tournament.GetGames())
				gameID := game.GetID()
				gameName := fmt.Sprintf("%s %d", trans("Game"), x+1)
				gameSlug := slug.Make(gameName) + "-" + util.RandomString(4)
				homeTeamID := uint(game.GetHomeTeam().GetID())
				awayTeamID := uint(game.GetAwayTeam().GetID())
				gameModel := models.Game{
					Name:         gameName,
					Slug:         controller.createUniqueGameSlug(gameSlug, 0),
					TournamentID: tournamentModel.ID,
					Teams:        gameTeams,
					Depth:        &depth,
				}

				if int(homeTeamID) != -1 {
					gameModel.HomeTeamID = &homeTeamID
				}
				if int(awayTeamID) != -1 {
					gameModel.AwayTeamID = &awayTeamID
				}

				res = controller.db.Save(&gameModel)
				if res.Error != nil {
					return slugString, res.Error
				}

				for _, pID := range game.GetParentIDs() {
					if gID, ok := insertedGames[pID]; ok {
						parentGameModel := models.Game{}
						parentGameModel.ID = uint(gID)
						res = controller.db.Where(parentGameModel).First(&parentGameModel)
						if res.Error != nil {
							return slugString, res.Error
						}
						parentGameModel.ChildID = &gameModel.ID
						res = controller.db.Save(&parentGameModel)
						if res.Error != nil {
							return slugString, res.Error
						}
					}
				}

				insertedGames[gameID] = int(gameModel.ID)
			}
		}
	} else {
		return slugString, errors.New("tournament type is invalid")
	}
	return slugString, nil
}

// getGameDepth gets the depth of the game in a tournament such as an elimination tournament. It is the same as counting how many games each team had to win in order to get to this game (a team which is by itself in a game automatically wins).
func getGameDepth(game tournify.GameInterface, games []tournify.GameInterface) int {
	ps := game.GetParentIDs()
	for _, id := range ps {
		for gameID, foundGame := range games {
			if gameID == id {
				return 1 + getGameDepth(foundGame, games)
			}
		}
	}
	return 0
}
