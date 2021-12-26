package models

import (
	"fmt"
	"github.com/tournify/tournify"
)

// Statistics is a default struct used as an example of how structs can be implemented for tournify
type Statistics struct {
	Tournament    tournify.TournamentInterface `json:"-"`
	Group         tournify.GroupInterface      `json:"-"`
	Team          Team                         `json:"team"`
	Played        int                          `json:"played"`
	Wins          int                          `json:"wins"`
	Losses        int                          `json:"losses"`
	Ties          int                          `json:"ties"`
	PointsFor     float64                      `json:"points_for"`
	PointsAgainst float64                      `json:"points_against"`
	Points        int                          `json:"points"`
}

// GetGroup returns the Group that the statistics were generated for, stats are directly related to a team and the group they are in.
func (t *Statistics) GetGroup() tournify.GroupInterface {
	return t.Group
}

// GetTeam returns the Team that the statistics were generated for, stats are directly related to a team and the group they are in.
func (t *Statistics) GetTeam() tournify.TeamInterface {
	return &t.Team
}

func (t *Statistics) GetTeamName() string {
	return t.Team.Name
}

// GetPlayed returns the number of games played
func (t *Statistics) GetPlayed() int {
	return t.Played
}

// GetWins returns the number of won games
func (t *Statistics) GetWins() int {
	return t.Wins
}

// GetLosses returns the number of lost games
func (t *Statistics) GetLosses() int {
	return t.Losses
}

// GetTies returns the number of games resulting in a tied game
func (t *Statistics) GetTies() int {
	return t.Ties
}

// GetPointsFor returns the number of goals or points that this team has made
func (t *Statistics) GetPointsFor() float64 {
	return t.PointsFor
}

// GetPointsAgainst returns the number of goals or points that other teams have made against this team
func (t *Statistics) GetPointsAgainst() float64 {
	return t.PointsAgainst
}

// GetDiff returns the difference of PointsFor and PointsAgainst
func (t *Statistics) GetDiff() float64 {
	return t.PointsFor - t.PointsAgainst
}

// GetPoints returns the number of points the team has based on wins, losses or ties
func (t *Statistics) GetPoints() int {
	return t.Points
}

// AddPoints adds the specified number of points to Points
func (t *Statistics) AddPoints(points int) {
	t.Points += points
}

// Print prints the stats as a single string
func (t *Statistics) Print() string {
	return fmt.Sprintf("%d\t%d\t%d\t%d\t%d\t%d\t%.0f/%.0f\t%.0f\t%d", t.GetGroup().GetID(), t.GetTeam().GetID(), t.GetPlayed(), t.GetWins(), t.GetTies(), t.GetLosses(), t.GetPointsFor(), t.GetPointsAgainst(), t.GetDiff(), t.GetPoints())
}
