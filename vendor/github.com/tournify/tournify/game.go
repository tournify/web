package tournify

import (
	"fmt"
)

// GameInterface defines the needed methods for games used in the library.
// A Game is a flexible entity and conforms to what you might typically find in Soccer where
// you have a home and away team and a score for each team but the interface also tries to
// allow for other types of games where the number of teams and scores is not limited to 2
type GameInterface interface {
	GetID() int
	GetParentIDs() []int
	GetHomeTeam() TeamInterface
	GetAwayTeam() TeamInterface
	SetHomeTeam(t TeamInterface)
	SetAwayTeam(t TeamInterface)
	GetHomeScore() ScoreInterface
	GetAwayScore() ScoreInterface
	GetTeams() []TeamInterface   // For games that can have any number of teams
	GetScores() []ScoreInterface // For games that can have any number of scores
	SetScore(homeScore float64, awayScore float64)
	Print() string
}

// Game is a default struct used as an example of how structs can be implemented for tournify
type Game struct {
	ID        int
	ParentIDs []int
	Scores    []ScoreInterface
	Teams     []TeamInterface
}

// GetParentIDs gets the ids of any games that caused this game to be generated, typically this is used in Elimination games
func (g *Game) GetParentIDs() []int {
	return g.ParentIDs
}

// SetScore sets home and away scores for home and away teams, this function is needed
// for games with a home and away team.
func (g *Game) SetScore(homeScore float64, awayScore float64) {
	if len(g.Scores) < 1 {
		g.Scores = append(g.Scores, &Score{}, &Score{})
	} else if len(g.Scores) < 2 {
		g.Scores = append(g.Scores, &Score{})
	}
	g.Scores[0].SetPoints(homeScore)
	g.Scores[1].SetPoints(awayScore)
}

// GetID returns the id of the game
func (g *Game) GetID() int {
	return g.ID
}

// GetHomeTeam returns the first team in the Teams slice
func (g *Game) GetHomeTeam() TeamInterface {
	if len(g.Teams) < 1 {
		g.Teams = append(g.Teams, &Team{})
	}
	return g.Teams[0]
}

// SetHomeTeam sets the first team of the game
func (g *Game) SetHomeTeam(t TeamInterface) {
	if len(g.Teams) < 1 {
		g.Teams = append(g.Teams, &Team{})
	}
	g.Teams[0] = t
}

// GetAwayTeam returns the second team in the Teams slice
func (g *Game) GetAwayTeam() TeamInterface {
	if len(g.Teams) < 1 {
		g.Teams = append(g.Teams, &Team{}, &Team{})
	} else if len(g.Teams) < 2 {
		g.Teams = append(g.Teams, &Team{})
	}
	return g.Teams[1]
}

// SetAwayTeam sets the second team of the game and adds a placeholder home team if it's not already there
func (g *Game) SetAwayTeam(t TeamInterface) {
	if len(g.Teams) < 1 {
		g.Teams = append(g.Teams, &Team{}, &Team{})
	} else if len(g.Teams) < 2 {
		g.Teams = append(g.Teams, &Team{})
	}
	g.Teams[1] = t
}

// GetHomeScore returns the first score in the Scores slice
func (g *Game) GetHomeScore() ScoreInterface {
	if len(g.Scores) < 1 {
		g.Scores = append(g.Scores, &Score{})
	}
	return g.Scores[0]
}

// GetAwayScore returns the second score in the Scores slice
func (g *Game) GetAwayScore() ScoreInterface {
	if len(g.Scores) < 1 {
		g.Scores = append(g.Scores, &Score{}, &Score{})
	} else if len(g.Scores) < 2 {
		g.Scores = append(g.Scores, &Score{})
	}
	return g.Scores[1]
}

// GetTeams returns a slice of Team
func (g *Game) GetTeams() []TeamInterface {
	return g.Teams
}

// GetScores returns a slice of Score
func (g *Game) GetScores() []ScoreInterface {
	return g.Scores
}

// Print writes game details to stdout
func (g *Game) Print() string {
	return fmt.Sprintf("Game ID: %d, HomeTeam: %d, AwayTeam: %d, HomeScore: %.2f, AwayScore: %.2f\n",
		g.GetID(),
		g.GetHomeTeam().GetID(),
		g.GetAwayTeam().GetID(),
		g.GetHomeScore().GetPoints(),
		g.GetAwayScore().GetPoints())
}
