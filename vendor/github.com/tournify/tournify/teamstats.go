package tournify

import (
	"errors"
	"fmt"
	"sort"
)

// TeamStatsInterface is used to show team statistics. Currently this is specifically made for
// group tournaments where there is a need to rank teams.
type TeamStatsInterface interface {
	GetGroup() GroupInterface
	GetTeam() TeamInterface
	GetPlayed() int
	GetWins() int
	GetLosses() int
	GetTies() int
	GetPointsFor() float64
	GetPointsAgainst() float64
	GetDiff() float64
	GetPoints() int
	AddPoints(points int)
	Print() string
}

// TeamStats is a default struct used as an example of how structs can be implemented for tournify
type TeamStats struct {
	Tournament    TournamentInterface
	Group         GroupInterface
	Team          TeamInterface
	Played        int
	Wins          int
	Losses        int
	Ties          int
	PointsFor     float64
	PointsAgainst float64
	Points        int
}

// GetGroup returns the Group that the statistics were generated for, stats are directly related to a team and the group they are in.
func (t *TeamStats) GetGroup() GroupInterface {
	return t.Group
}

// GetTeam returns the Team that the statistics were generated for, stats are directly related to a team and the group they are in.
func (t *TeamStats) GetTeam() TeamInterface {
	return t.Team
}

// GetPlayed returns the number of games played
func (t *TeamStats) GetPlayed() int {
	return t.Played
}

// GetWins returns the number of won games
func (t *TeamStats) GetWins() int {
	return t.Wins
}

// GetLosses returns the number of lost games
func (t *TeamStats) GetLosses() int {
	return t.Losses
}

// GetTies returns the number of games resulting in a tied game
func (t *TeamStats) GetTies() int {
	return t.Ties
}

// GetPointsFor returns the number of goals or points that this team has made
func (t *TeamStats) GetPointsFor() float64 {
	return t.PointsFor
}

// GetPointsAgainst returns the number of goals or points that other teams have made against this team
func (t *TeamStats) GetPointsAgainst() float64 {
	return t.PointsAgainst
}

// GetDiff returns the difference of PointsFor and PointsAgainst
func (t *TeamStats) GetDiff() float64 {
	return t.PointsFor - t.PointsAgainst
}

// GetPoints returns the number of points the team has based on wins, losses or ties
func (t *TeamStats) GetPoints() int {
	return t.Points
}

// AddPoints adds the specified number of points to Points
func (t *TeamStats) AddPoints(points int) {
	t.Points += points
}

// Print prints the stats as a single string
func (t *TeamStats) Print() string {
	return fmt.Sprintf("%d\t%d\t%d\t%d\t%d\t%d\t%.0f/%.0f\t%.0f\t%d", t.GetGroup().GetID(), t.GetTeam().GetID(), t.GetPlayed(), t.GetWins(), t.GetTies(), t.GetLosses(), t.GetPointsFor(), t.GetPointsAgainst(), t.GetDiff(), t.GetPoints())
}

// PrintGroupTournamentStats takes an array of team stats and returns them as a single string with a new line for each team
func PrintGroupTournamentStats(teamStats []TeamStatsInterface) string {
	res := "Group\tTeam\tPlayed\tWins\tTies\tLosses\t+/-\tDiff\tPoints\n"
	for _, ts := range teamStats {
		res += ts.Print()
		res += "\n"
	}
	return res
}

// GetGroupTournamentStats takes 4 inputs. The first input is the tournament itself.
// The other three input defines how many points a team should get for a win, loss or tie. The standard is 3, 0, 1 but
// it can vary depending on the tournament.
func GetGroupTournamentStats(t TournamentInterface, winPoints int, lossPoints int, tiePoints int) ([]TeamStatsInterface, error) {
	if t.GetType() != int(TournamentTypeGroup) {
		return nil, errors.New("can not get stats for tournament type which is not TournamentTypeGroup")
	}
	var stats []TeamStatsInterface

	for _, group := range t.GetGroups() {
		groupStats := GetGroupStats(group, winPoints, lossPoints, tiePoints)
		stats = append(stats, groupStats...)
	}
	return stats, nil
}

// GetGroupStats returns the current stats for all the teams in a single group
func GetGroupStats(group GroupInterface, winPoints int, lossPoints int, tiePoints int) []TeamStatsInterface {
	var groupStats []TeamStatsInterface
	teamStats := map[int]*TeamStats{}

	for _, team := range *group.GetTeams() {
		teamStats[team.GetID()] = &TeamStats{
			Group: group,
			Team:  team,
		}
	}

	for _, game := range *group.GetGames() {
		if _, ok := teamStats[game.GetHomeTeam().GetID()]; !ok {
			teamStats[game.GetHomeTeam().GetID()] = &TeamStats{
				Group: group,
				Team:  game.GetHomeTeam(),
			}
		}
		// Calculate stats for the home team in every game
		teamStats[game.GetHomeTeam().GetID()].PointsFor += game.GetHomeScore().GetPoints()
		teamStats[game.GetHomeTeam().GetID()].PointsAgainst += game.GetAwayScore().GetPoints()
		if game.GetHomeScore().GetPoints() > game.GetAwayScore().GetPoints() {
			teamStats[game.GetHomeTeam().GetID()].Wins++
		} else if game.GetHomeScore().GetPoints() == game.GetAwayScore().GetPoints() {
			teamStats[game.GetHomeTeam().GetID()].Ties++
		} else {
			teamStats[game.GetHomeTeam().GetID()].Losses++
		}

		teamStats[game.GetHomeTeam().GetID()].Played++

		// Calculate stats for the away team in every game
		if _, ok := teamStats[game.GetAwayTeam().GetID()]; !ok {
			teamStats[game.GetAwayTeam().GetID()] = &TeamStats{
				Group: group,
				Team:  game.GetAwayTeam(),
			}
		}
		teamStats[game.GetAwayTeam().GetID()].PointsFor += game.GetAwayScore().GetPoints()
		teamStats[game.GetAwayTeam().GetID()].PointsAgainst += game.GetHomeScore().GetPoints()
		if game.GetHomeScore().GetPoints() < game.GetAwayScore().GetPoints() {
			teamStats[game.GetAwayTeam().GetID()].Wins++
		} else if game.GetHomeScore().GetPoints() == game.GetAwayScore().GetPoints() {
			teamStats[game.GetAwayTeam().GetID()].Ties++
		} else {
			teamStats[game.GetAwayTeam().GetID()].Losses++
		}
		teamStats[game.GetAwayTeam().GetID()].Played++
	}

	for _, t := range teamStats {
		t.Points = t.Wins * winPoints
		t.Points += t.Losses * lossPoints
		t.Points += t.Ties * tiePoints
		groupStats = append(groupStats, t)
	}

	groupStats = SortTournamentStats(groupStats)
	return groupStats
}

// SortTournamentStats sorts the statistics by points, diff and finally scored goals against other teams
func SortTournamentStats(stats []TeamStatsInterface) []TeamStatsInterface {
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].GetPoints() > stats[j].GetPoints() {
			return true
		} else if stats[i].GetPoints() < stats[j].GetPoints() {
			return false
		} else {
			if stats[i].GetDiff() > stats[j].GetDiff() {
				return true
			} else if stats[i].GetDiff() < stats[j].GetDiff() {
				return false
			} else {
				if stats[i].GetPointsFor() > stats[j].GetPointsFor() {
					return true
				} else if stats[i].GetPointsFor() < stats[j].GetPointsFor() {
					return false
				}
			}
		}
		return true
	})
	return stats
}
