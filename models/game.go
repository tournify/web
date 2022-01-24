package models

import (
	"fmt"
	"github.com/tournify/tournify"
	"gorm.io/gorm"
)

type Game struct {
	gorm.Model
	Name         string     `json:"name"`
	Slug         string     `gorm:"uniqueIndex;size:256;" json:"slug"`
	Keywords     string     `json:"-"`
	Description  string     `json:"-"`
	HomeTeamID   *uint      `json:"home_team"`
	AwayTeamID   *uint      `json:"away_team"`
	Teams        []Team     `gorm:"many2many:game_teams;" json:"teams"`
	Tournament   Tournament `json:"-"`
	TournamentID uint       `json:"-"`
	Group        Group      `json:"-"`
	GroupID      *uint      `json:"-"`
	ChildID      *uint      `json:"-"`
	Parents      []Game     `gorm:"foreignKey:ChildID" json:"-"`
	Scores       []Score    `json:"scores"`
	Depth        *int       `json:"depth,omitempty"`
}

// GetParentIDs gets the ids of any games that caused this game to be generated, typically this is used in Elimination games
func (g *Game) GetParentIDs() []int {
	var res []int
	for _, pg := range g.Parents {
		res = append(res, int(pg.ID))
	}
	return res
}

// SetScore sets home and away scores for home and away teams, this function is needed
// for games with a home and away team.
func (g *Game) SetScore(homeScore float64, awayScore float64) {
	if len(g.Scores) < 1 {
		g.Scores = append(g.Scores, Score{
			TeamID: uint(g.GetHomeTeam().GetID()),
		}, Score{
			TeamID: uint(g.GetAwayTeam().GetID()),
		})
	} else if len(g.Scores) < 2 {
		g.Scores = append(g.Scores, Score{
			TeamID: uint(g.GetAwayTeam().GetID()),
		})
	}
	g.Scores[0].SetPoints(homeScore)
	g.Scores[1].SetPoints(awayScore)
}

// GetID returns the id of the game
func (g *Game) GetID() int {
	return int(g.ID)
}

// GetHomeTeam returns the first team in the Teams slice
func (g *Game) GetHomeTeam() tournify.TeamInterface {
	if len(g.Teams) < 1 {
		gameID := -1
		g.Teams = append(g.Teams, Team{
			Model: gorm.Model{
				ID: uint(gameID),
			},
		})
	}
	if g.HomeTeamID != nil {
		for i, t := range g.Teams {
			if int(t.ID) == int(*g.HomeTeamID) {
				return &g.Teams[i]
			}
		}
	}
	return &g.Teams[0]
}

func (g *Game) GetHomeTeamName() string {
	if g.HomeTeamID != nil {
		for i, t := range g.Teams {
			if int(t.ID) == int(*g.HomeTeamID) {
				return g.Teams[i].Name
			}
		}
	}
	if len(g.Teams) >= 1 {
		return g.Teams[0].Name
	}
	return ""
}

// SetHomeTeam sets the first team of the game
func (g *Game) SetHomeTeam(t tournify.TeamInterface) {
	team := t.(*Team)
	if len(g.Teams) < 1 {
		g.Teams = append(g.Teams, *team)
	}
	if g.HomeTeamID != nil {
		for i := range g.Teams {
			if g.Teams[i].GetID() == int(*g.HomeTeamID) {
				g.Teams[i] = *team
			}
		}
	} else {
		g.Teams[0] = *team
	}
	g.HomeTeamID = &team.ID
}

// GetAwayTeam returns the second team in the Teams slice
func (g *Game) GetAwayTeam() tournify.TeamInterface {
	if g.AwayTeamID != nil {
		for i := range g.Teams {
			if g.Teams[i].GetID() == int(*g.AwayTeamID) {
				return &g.Teams[i]
			}
		}
	}
	if len(g.Teams) < 1 {
		gameID := -1
		g.Teams = append(g.Teams, Team{
			Model: gorm.Model{
				ID: uint(gameID),
			},
		}, Team{
			Model: gorm.Model{
				ID: uint(gameID),
			},
		})
	} else if len(g.Teams) < 2 {
		gameID := -1
		g.Teams = append(g.Teams, Team{
			Model: gorm.Model{
				ID: uint(gameID),
			},
		})
	}
	if g.AwayTeamID != nil {
		for i, t := range g.Teams {
			if int(t.ID) == int(*g.AwayTeamID) {
				return &g.Teams[i]
			}
		}
	}
	return &g.Teams[1]
}

func (g *Game) GetAwayTeamName() string {
	if g.AwayTeamID != nil {
		for i, t := range g.Teams {
			if int(t.ID) == int(*g.AwayTeamID) {
				return g.Teams[i].Name
			}
		}
	}
	if len(g.Teams) >= 2 {
		return g.Teams[1].Name
	}
	return ""
}

// SetAwayTeam sets the second team of the game and adds a placeholder home team if it's not already there
func (g *Game) SetAwayTeam(t tournify.TeamInterface) {
	team := t.(*Team)
	if len(g.Teams) < 1 {
		g.Teams = append(g.Teams, *team)
	} else if g.AwayTeamID != nil {
		for i := range g.Teams {
			if g.Teams[i].GetID() == int(*g.AwayTeamID) {
				g.Teams[i] = *team
			}
		}
	} else if len(g.Teams) == 1 {
		if int(g.Teams[0].ID) == -1 {
			g.Teams[0] = *team
		} else {
			g.Teams = append(g.Teams, *team)
		}
	} else if len(g.Teams) >= 2 {
		g.Teams[1] = *team
	}
	g.AwayTeamID = &team.ID
}

// GetHomeScore returns the first score in the Scores slice
func (g *Game) GetHomeScore() tournify.ScoreInterface {
	if len(g.Scores) < 1 {
		g.Scores = append(g.Scores, Score{})
	}
	if g.HomeTeamID != nil {
		for i, s := range g.Scores {
			if int(s.TeamID) == int(*g.HomeTeamID) {
				return &g.Scores[i]
			}
		}
	}
	return &g.Scores[0]
}

// GetAwayScore returns the second score in the Scores slice
func (g *Game) GetAwayScore() tournify.ScoreInterface {
	if len(g.Scores) < 1 {
		g.Scores = append(g.Scores, Score{}, Score{})
	} else if len(g.Scores) < 2 {
		g.Scores = append(g.Scores, Score{})
	}
	if g.AwayTeamID != nil {
		for i, s := range g.Scores {
			if int(s.TeamID) == int(*g.AwayTeamID) {
				return &g.Scores[i]
			}
		}
	}
	return &g.Scores[1]
}

// GetTeams returns a slice of Team
func (g *Game) GetTeams() []tournify.TeamInterface {
	var teams []tournify.TeamInterface
	for _, t := range g.Teams {
		teams = append(teams, &t)
	}
	return teams
}

// GetScores returns a slice of Score
func (g *Game) GetScores() []tournify.ScoreInterface {
	var scores []tournify.ScoreInterface
	for _, s := range g.Scores {
		scores = append(scores, &s)
	}
	return scores
}

// Print writes game details to stdout
func (g *Game) Print() string {
	return fmt.Sprintf("Game ID: %d, HomeTeamID: %d, AwayTeamID: %d, HomeScore: %.2f, AwayScore: %.2f\n",
		g.GetID(),
		g.GetHomeTeam().GetID(),
		g.GetAwayTeam().GetID(),
		g.GetHomeScore().GetPoints(),
		g.GetAwayScore().GetPoints())
}
