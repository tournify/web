package tournify

import (
	"fmt"
)

// TeamInterface defines the methods for teams. Teams are used to create tournaments and generate games.
// Teams can have games
type TeamInterface interface {
	GetID() int
	GetEliminatedCount() int
	Print() string
}

// Team is a default struct used as an example of how structs can be implemented for tournify
type Team struct {
	ID         int
	Eliminated int // Increment by 1 every time this team is elimnated
}

// GetID returns the id of the score
func (t *Team) GetID() int {
	return t.ID
}

// Print writes team details to stdout
func (t *Team) Print() string {
	return fmt.Sprintf("Team ID: %d\n", t.GetID())
}

// GetEliminatedCount gets the number of times the team has been eliminated in a multiple elimination tournament
func (t *Team) GetEliminatedCount() int {
	return t.Eliminated
}
