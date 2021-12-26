package tournify

// TournamentType defines the type of tournament
// TODO is this a good way of defining tournament types?
//  How can this be extended by others who want to build on this code without modifying the library?
type TournamentType int

const (
	// TournamentTypeGroup is for group tournaments
	TournamentTypeGroup TournamentType = 0
	// TournamentTypeSeries is similar to group tournaments but only has one group, usually this would be used when playing in a soccer league for example
	TournamentTypeSeries TournamentType = 1
	// TournamentTypeElimination is for elimination or knockout tournaments
	TournamentTypeElimination TournamentType = 2
	// TournamentTypeDoubleElimination is the same as TournamentTypeElimination
	// but teams to get knocked out early get a second chance to come back and win
	TournamentTypeDoubleElimination TournamentType = 3
)

func (tournamentType TournamentType) String() string {
	names := [...]string{"Group", "Series", "Elimination", "Double Elimination"}

	if tournamentType < TournamentTypeGroup || tournamentType > TournamentTypeDoubleElimination {
		return "Unknown"
	}

	return names[tournamentType]
}
