package tournify

import (
	"errors"
	"fmt"
	"math"
)

// TournamentInterface defines the methods needed to handle tournaments.
type TournamentInterface interface {
	GetType() int
	GetTypeString() string
	GetTeams() []TeamInterface
	GetEliminatedTeams() []TeamInterface // For elimination style tournaments
	GetRemainingTeams() []TeamInterface  // For elimination style tournaments
	GetGroups() []GroupInterface
	GetGames() []GameInterface
	AppendGame(game GameInterface)
	SetGame(game GameInterface) error
	SetGameScore(game GameInterface, homeScore float64, awayScore float64) error
	Print() string
}

// Tournament is a default struct used as an example of how structs can be implemented for tournify
type Tournament struct {
	Type   TournamentType // Is it elimination or group or ladder or poker? What is a type?
	Teams  []TeamInterface
	Groups []GroupInterface
	Games  []GameInterface
}

// GetType returns the type of tournament as an int
func (t *Tournament) GetType() int {
	return int(t.Type)
}

// GetTypeString returns the type of tournament as a string
func (t *Tournament) GetTypeString() string {
	return t.Type.String()
}

// GetTeams returns the team slice
func (t *Tournament) GetTeams() []TeamInterface {
	return t.Teams
}

// GetGroups returns the group slice
func (t *Tournament) GetGroups() []GroupInterface {
	return t.Groups
}

// GetGames returns the game slice
func (t *Tournament) GetGames() []GameInterface {
	return t.Games
}

// SetGames sets the tournaments games slice
func (t *Tournament) SetGames(games []GameInterface) {
	for _, g := range games {
		err := t.SetGame(g)
		if err != nil {
			return
		}
	}
}

// GetEliminatedTeams gets all teams that have been eliminated at least one time in an elimination tournament
func (t *Tournament) GetEliminatedTeams() []TeamInterface {
	var elimnatedTeams []TeamInterface
	for _, team := range t.GetTeams() {
		if team.GetEliminatedCount() > 0 && t.GetType() == int(TournamentTypeElimination) {
			elimnatedTeams = append(elimnatedTeams, team)
		}
		if team.GetEliminatedCount() > 1 && t.GetType() == int(TournamentTypeDoubleElimination) {
			elimnatedTeams = append(elimnatedTeams, team)
		}
	}
	return elimnatedTeams
}

// GetRemainingTeams gets all teams that have not been eliminated in an elimination tournament
func (t *Tournament) GetRemainingTeams() []TeamInterface {
	var remainingTeams []TeamInterface
	for _, team := range t.GetTeams() {
		if team.GetEliminatedCount() < 1 && t.GetType() == int(TournamentTypeElimination) {
			remainingTeams = append(remainingTeams, team)
		}
		if team.GetEliminatedCount() < 2 && t.GetType() == int(TournamentTypeDoubleElimination) {
			remainingTeams = append(remainingTeams, team)
		}
	}
	return remainingTeams
}

// AppendGame appends a game to the tournament game slice
func (t *Tournament) AppendGame(game GameInterface) {
	t.Games = append(t.Games, game)
}

// SetGame overwrites any game with the same id
func (t *Tournament) SetGame(game GameInterface) error {
	if TournamentType(t.GetType()) == TournamentTypeElimination {
		return t.setEliminationGame(game)
	} else {
		for i, g := range t.Games {
			if g.GetID() == game.GetID() {
				t.Games[i] = game
				return nil
			}
		}
	}
	return errors.New("could not set game, no matching game ID found")
}

func (t *Tournament) setEliminationGame(game GameInterface) error {
	for i, g := range t.Games {
		if g.GetID() == game.GetID() {
			// check if this changes winner or loser, if it does we need to change the child of this game
			if WinnerChanged(g, game) {
				// Find the game which has this game as a parent
				for _, tg := range t.Games {
					for _, parentID := range tg.GetParentIDs() {
						if parentID == g.GetID() {
							winTeam := GetWinnerTeam(game)
							loseTeam := GetLoserTeam(game)
							if tg.GetHomeTeam().GetID() == -1 || loseTeam.GetID() == tg.GetHomeTeam().GetID() {
								tg.SetHomeTeam(winTeam)
							} else if tg.GetAwayTeam().GetID() == -1 || loseTeam.GetID() == tg.GetHomeTeam().GetID() {
								tg.SetAwayTeam(winTeam)
							}
							t.MakeTeamEliminated(loseTeam)
							t.MakeTeamRemain(winTeam)
						}
					}
				}
			}
			t.Games[i] = game
			return nil
		}
	}
	return errors.New("could not set game, no matching game ID found")
}

func (t *Tournament) SetGameScore(game GameInterface, homeScore float64, awayScore float64) error {
	if t.GetType() == int(TournamentTypeElimination) {
		oldGame := Game{
			ID:        game.GetID(),
			ParentIDs: game.GetParentIDs(),
		}
		for _, s := range game.GetScores() {
			oldGame.Scores = append(oldGame.Scores, &Score{
				ID:     s.GetID(),
				Points: s.GetPoints(),
			})
		}
		for _, team := range game.GetTeams() {
			oldGame.Teams = append(oldGame.Teams, &Team{
				ID:         team.GetID(),
				Eliminated: team.GetEliminatedCount(),
			})
		}
		game.SetScore(homeScore, awayScore)
		if WinnerChanged(&oldGame, game) {
			for _, tg := range t.Games {
				for _, parentID := range tg.GetParentIDs() {
					if parentID == game.GetID() {
						winTeam := GetWinnerTeam(game)
						loseTeam := GetLoserTeam(game)
						// Win team is nil in a tie game
						if winTeam != nil {
							if tg.GetHomeTeam().GetID() == -1 || loseTeam.GetID() == tg.GetHomeTeam().GetID() {
								tg.SetHomeTeam(winTeam)
							} else if tg.GetAwayTeam().GetID() == -1 || loseTeam.GetID() == tg.GetHomeTeam().GetID() {
								tg.SetAwayTeam(winTeam)
							}
							t.MakeTeamEliminated(loseTeam)
							t.MakeTeamRemain(winTeam)
							err := t.SetGame(game)
							if err != nil {
								return err
							}
							return nil
						}
					}
				}
			}
		}
	} else {
		game.SetScore(homeScore, awayScore)
		err := t.SetGame(game)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Tournament) MakeTeamEliminated(team TeamInterface) {
	if TournamentType(t.GetType()) == TournamentTypeElimination {
		for i, tt := range t.Teams {
			if tt.GetID() == team.GetID() {
				t.Teams[i].SetEliminatedCount(1)
				return
			}
		}
	}
}

func (t *Tournament) MakeTeamRemain(team TeamInterface) {
	if t.Type == TournamentTypeElimination {
		for i, tt := range t.Teams {
			if tt.GetID() == team.GetID() {
				t.Teams[i].SetEliminatedCount(0)
				return
			}
		}
	}
}

// IsDepthFull checks if the expected numbers of games have been filled for a depth in an elimination tournament
func (t *Tournament) IsDepthFull(depth int) bool {
	if depth <= 0 {
		return true
	}
	c := len(t.GetGamesAtDepth(0))
	// Fix for uneven count
	if c%2 != 0 {
		c = +1
	}
	c = c / (2 * depth)
	// It has to be less than or equal here
	// start with 10 games, depth is 2 therefore c is 2.5 and games at depth when full is 3 but 2 games is not full
	// start with 8 games, depth is 2 therefore c is 2 and games at depth when full is 2
	if c <= len(t.GetGamesAtDepth(depth)) {
		return true
	}
	return false
}

// GetGamesAtDepth takes an int for depth and returns any games at the depth as a slice
func (t *Tournament) GetGamesAtDepth(depth int) (games []GameInterface) {
	for _, g := range t.GetGames() {
		if t.GetGameDepth(g) == depth {
			games = append(games, g)
		}
	}
	return games
}

// GetGameDepth gets the depth of the game in a tournament such as an elimination tournament. It is the same as counting how many games each team had to win in order to get to this game (a team which is by itself in a game automatically wins).
func (t *Tournament) GetGameDepth(game GameInterface) int {
	ps := game.GetParentIDs()
	for _, id := range ps {
		return 1 + t.GetGameDepth(t.GetGameByID(id))
	}
	return 0
}

// GetGameByID takes an int and returns a game with that id if it exists in the tournament
func (t *Tournament) GetGameByID(id int) GameInterface {
	for _, g := range t.GetGames() {
		if g.GetID() == id {
			return g
		}
	}
	return nil
}

// GetGameFirstAncestorID gets the lowest game id with a depth of 0 which this game is a descendant of
func (t *Tournament) GetGameFirstAncestorID(game GameInterface) int {
	ps := game.GetParentIDs()
	if len(ps) > 0 {
		lowestParent := 0
		for _, pid := range ps {
			if pid > lowestParent {
				lowestParent = pid
			}
		}
		for _, g := range t.GetGames() {
			if g.GetID() == lowestParent {
				return t.GetGameFirstAncestorID(g)
			}
		}
	}
	return game.GetID()
}

// GetGameLastDescendantID gets the id of the last game that has been generated off of the provided game
func (t *Tournament) GetGameLastDescendantID(game GameInterface) int {
	for _, g := range t.GetGames() {
		for _, id := range g.GetParentIDs() {
			if id == game.GetID() {
				return t.GetGameLastDescendantID(g)
			}
		}
	}
	return game.GetID()
}

// Print writes the full tournament details to a string
func (t *Tournament) Print() string {
	p := fmt.Sprintf("TournamentType: %s\n", t.GetTypeString())
	if t.GetType() == 0 {
		p += fmt.Sprintf("\nGroups\n")
		for _, group := range t.GetGroups() {
			p += group.Print()
		}
	} else {
		p += fmt.Sprintf("\nTeams\n")
		for _, team := range t.GetTeams() {
			p += team.Print()
		}
	}
	p += fmt.Sprintf("\nGames\n")
	for _, games := range t.GetGames() {
		p += games.Print()
	}
	return p
}

// CreateTournament creates a tournament with the simplest input. It is recommended to create a slice with
// specific use via CreateTournamentFromTeams as this method will generate it's own Teams as a sort of placeholder.
func CreateTournament(teamCount int, meetCount int, groupCount int, tournamentType int) TournamentInterface {
	var teams []TeamInterface

	for i := 0; i < teamCount; i++ {
		teams = append(teams, &Team{ID: i})
	}

	return CreateTournamentFromTeams(teams, meetCount, groupCount, tournamentType)
}

// CreateTournamentFromTeams takes a slice of teams and generates a tournament of the specified type
func CreateTournamentFromTeams(teams []TeamInterface, meetCount int, groupCount int, tournamentType int) TournamentInterface {
	if TournamentType(tournamentType) == TournamentTypeGroup {
		if groupCount < 1 {
			return nil
		}
		if meetCount < 1 {
			return nil
		}
		return CreateGroupTournamentFromTeams(teams, groupCount, meetCount)
	} else if TournamentType(tournamentType) == TournamentTypeSeries {
		// TODO this should return an tournament of type series
		return CreateGroupTournamentFromTeams(teams, 1, meetCount)
	} else if TournamentType(tournamentType) == TournamentTypeElimination {
		return CreateEliminationTournamentFromTeams(teams)
	}
	return nil
}

// CreateEliminationTournamentFromTeams takes a slice of teams and generates a elimination tournament
// The ID used for games are very important for elimination tournaments as it is used to determine the home or away team in later games
func CreateEliminationTournamentFromTeams(teams []TeamInterface) TournamentInterface {
	// Create the initial games of the elimination tournament
	var games []GameInterface
	gameID := 0
	for i := 0; i < len(teams); i += 2 {
		game := Game{ID: gameID, Teams: []TeamInterface{teams[i]}}
		if i+1 < len(teams) {
			game.SetAwayTeam(teams[i+1])
		}
		gameID++
		games = append(games, &game)
	}
	games = createDescendantGames(games, gameID)
	// Return a tournament
	return &Tournament{Games: games, Teams: teams, Type: TournamentTypeElimination}
}

func createDescendantGames(gs []GameInterface, gameID int) []GameInterface {
	var newGames []GameInterface
	for i := 0; i < len(gs); i += 2 {
		game := Game{ID: gameID, ParentIDs: []int{
			gs[i].GetID(),
		}}
		if i+1 < len(gs) {
			game.ParentIDs = append(game.ParentIDs, gs[i+1].GetID())
		}
		gameID++
		newGames = append(newGames, &game)
	}
	if len(newGames) <= 1 {
		return append(gs, newGames...)
	}
	newGames = createDescendantGames(newGames, gameID)
	return append(gs, newGames...)
}

// CreateGroupTournamentFromTeams takes a slice of teams and generates a group tournament
func CreateGroupTournamentFromTeams(teams []TeamInterface, groupCount int, meetCount int) TournamentInterface {
	// TODO implement better error handling
	if groupCount < 1 || meetCount < 1 {
		return nil
	}

	groups := []GroupInterface{&Group{ID: 0}}
	teamsPerGroup := len(teams) / groupCount

	for i := 1; i < groupCount; i++ {
		groups = append(groups, &Group{ID: i})
	}

	groupIndex := 0
	for i, team := range teams {
		adjGI := groupIndex + 1
		if i >= teamsPerGroup*adjGI && adjGI < groupCount {
			groupIndex++
		}
		groups[groupIndex].AppendTeam(team)
	}

	return CreateGroupTournamentFromGroups(groups, meetCount)
}

// CreateGroupTournamentFromGroups takes a slice of groups that contain teams and returns a group tournament
// TODO simplify and break down this function in to smaller chunks?
// TODO this method currently uses cross matching for games but other types of matching could be supported
func CreateGroupTournamentFromGroups(groups []GroupInterface, meetCount int) TournamentInterface {
	// Works best for an even amount of teams in every group
	var games []GameInterface
	var teams []TeamInterface
	gameIndex := 0
	for gi, group := range groups {
		var tempID int
		uneven := false

		teams = append(teams, *group.GetTeams()...)
		gTeams := *group.GetTeams()

		// If there is an uneven amount of teams we need to add a temporary team which is later removed
		if len(gTeams)%2 != 0 {
			tempID = generateTempID(gTeams, -1)
			tempTeam := Team{ID: tempID}
			gTeams = append(gTeams, &tempTeam)
			uneven = true
		}

		// Loop through meet count
		for mi := 0; mi < meetCount; mi++ {
			// TODO game calculation is wrong when there is an uneven number of teams per group
			if len(gTeams) > 1 {
				halfCountHiger := DivideRoundUp(len(gTeams), 2)
				halfCountLower := DivideRoundDown(len(gTeams), 2)
				homeTeams := make([]TeamInterface, halfCountHiger)
				awayTeams := make([]TeamInterface, halfCountLower)
				// Everyone meets everyone once
				// We begin by taking our slice of teams like 0,1,2,3, and splitting it into home and away teams
				// if meet index is even
				if mi%2 == 0 {
					// The first half of the team slice become the home teams
					copy(homeTeams, gTeams[0:halfCountHiger])
					// The second half of the team slice become the away teams
					copy(awayTeams, gTeams[halfCountHiger:])
					// if meet index is odd
				} else {
					copy(awayTeams, gTeams[0:halfCountHiger])
					copy(homeTeams, gTeams[halfCountLower:])
				}

				awayTeams = reverseSlice(awayTeams)

				for i := 0; i < len(gTeams)-1; i++ {
					// Now we have home teams of 0,1 and away teams of 2,3
					// This means 0 will meet 2 and 1 will meet 3
					for hi, hteam := range homeTeams {
						game := Game{ID: gameIndex, Teams: []TeamInterface{hteam, awayTeams[hi]}}
						groups[gi].AppendGame(&game)
						games = append(games, &game)
						gameIndex++
					}
					homeTeams, awayTeams = rotateTeamsForCrossMatching(homeTeams, awayTeams)

				}
			}
		}
		if uneven {
			var removedGames []GameInterface
			games, removedGames = removeTempGames(games, removedGames, tempID)
			for _, removedGame := range removedGames {
				groups[gi].RemoveGame(removedGame)
			}
		}
	}
	return &Tournament{Groups: groups, Games: games, Teams: teams, Type: TournamentTypeGroup}
}

func removeTempGames(games []GameInterface, removedGames []GameInterface, tempID int) ([]GameInterface, []GameInterface) {
	for i := 0; i < len(games); i++ {
		if games[i].GetHomeTeam().GetID() == tempID || games[i].GetAwayTeam().GetID() == tempID {
			removedGames = append(removedGames, games[i])
			tmpGames := append(games[:i], games[i+1:]...)
			return removeTempGames(tmpGames, removedGames, tempID)
		}
	}
	return games, removedGames
}

func generateTempID(teams []TeamInterface, tempID int) int {
	for _, t := range teams {
		if t.GetID() == tempID {
			return generateTempID(teams, tempID-1)
		}
	}
	return tempID
}

func reverseSlice(a []TeamInterface) []TeamInterface {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
	return a
}

func rotateTeamsForCrossMatching(homeTeams []TeamInterface, awayTeams []TeamInterface) ([]TeamInterface, []TeamInterface) {
	var x, y, z TeamInterface
	// We keep the first home team in the same position and rotate all others
	// HT = Home Teams, AT = Away Teams
	// for HT 0,1 and AT 2,3. 0 is kept in place while 1 remains in the home team array
	x, homeTeams = homeTeams[0], homeTeams[1:]
	// Take the first away team
	// 2 is taken out of AT, 3 remains in AT
	z, awayTeams = awayTeams[0], awayTeams[1:]
	// and append to end of home teams
	// HT is now 1,2
	homeTeams = append(homeTeams, z)
	// Take the first home team
	// 1 is taken out of HT, HT is now 2
	y, homeTeams = homeTeams[0], homeTeams[1:]
	// and append it to the end of away teams
	// 1 is added to end of AT, AT is now 3,1
	awayTeams = append(awayTeams, y)
	// Put the first home team back in first position of home array
	// HT is now 0,2
	homeTeams = append([]TeamInterface{x}, homeTeams...)
	return homeTeams, awayTeams
}

// NumberOfGamesForGroupTournament Calculates the number of games in a group tournament based on number of teams, groups and unique encounters. This is unreliable with uneven numbers.
func NumberOfGamesForGroupTournament(teamCount int, groupCount int, meetCount int) int {
	tpg := float64(teamCount) / float64(groupCount)
	games := tpg * (tpg - 1) / 2
	res := int(games * float64(meetCount*groupCount))
	if math.Mod(float64(teamCount), float64(groupCount)) != 0 {
		res += int(math.Mod(float64(teamCount), float64(groupCount))) * meetCount
	}
	return res
}

// NumberOfGamesForEliminationTournament Calculates the number of games in a elimination tournament based on the number of teams
func NumberOfGamesForEliminationTournament(teamCount int) int {
	if teamCount <= 1 {
		return 0
	}
	// If the team count starts as uneven we will only have 1 extra game
	if teamCount%2 != 0 {
		teamCount = teamCount - 1
		return 1 + NumberOfGamesForEliminationTournament(teamCount/2+1) + teamCount/2
	}
	return NumberOfGamesForEliminationTournament(teamCount/2) + teamCount/2
}

// DivideRoundUp takes two ints, divides them and rounds the result up to the nearest int
func DivideRoundUp(a int, b int) int {
	return int(math.Ceil(float64(a) / float64(b)))
}

// DivideRoundDown takes two ints, divides them and rounds the result up to the nearest int
func DivideRoundDown(a int, b int) int {
	return int(math.Floor(float64(a) / float64(b)))
}
