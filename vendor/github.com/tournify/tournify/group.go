package tournify

import (
	"fmt"
)

// GroupInterface defines the interface of tournament groups used for group tournaments
type GroupInterface interface {
	GetID() int
	GetTeams() *[]TeamInterface
	GetGames() *[]GameInterface
	AppendGames(games []GameInterface)
	AppendGame(game GameInterface)
	RemoveGame(game GameInterface)
	AppendTeams(teams []TeamInterface)
	AppendTeam(team TeamInterface)
	Print() string
}

// Group is for group tournaments only
type Group struct {
	ID    int
	Teams []TeamInterface
	Games []GameInterface
}

// GetID returns the id of the group
func (t *Group) GetID() int {
	return t.ID
}

// GetTeams returns a slice of teams belonging to the group
func (t *Group) GetTeams() *[]TeamInterface {
	return &t.Teams
}

// GetGames returns the slice of games belonging to the group
func (t *Group) GetGames() *[]GameInterface {
	return &t.Games
}

// AppendGames adds a slice of games to the Games slice
func (t *Group) AppendGames(games []GameInterface) {
	t.Games = append(t.Games, games...)
}

// AppendGame takes a single game and appends it to the Games slice
func (t *Group) AppendGame(game GameInterface) {
	t.Games = append(t.Games, game)
}

// RemoveGame takes a single game and removes it to the Games slice
func (t *Group) RemoveGame(game GameInterface) {
	for i := 0; i < len(t.Games); i++ {
		if t.Games[i].GetID() == game.GetID() {
			t.Games = append(t.Games[:i], t.Games[i+1:]...)
		}
	}
}

// AppendTeams adds a slice of teams to the Teams slice
func (t *Group) AppendTeams(teams []TeamInterface) {
	t.Teams = append(t.Teams, teams...)
}

// AppendTeam takes a single team and appends it to the Teams slice
func (t *Group) AppendTeam(team TeamInterface) {
	t.Teams = append(t.Teams, team)
}

// Print writes group details to stdout
func (t *Group) Print() string {
	p := fmt.Sprintf("Group ID: %d\n", t.GetID())
	for _, team := range *t.GetTeams() {
		p += team.Print()
	}
	p += fmt.Sprintf("\n")
	return p
}
