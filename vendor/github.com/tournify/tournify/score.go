package tournify

// ScoreInterface defines how scores should be defined. Scores hold points defined as a float.
// Using a float64 should allow this library to be used for any type of game.
type ScoreInterface interface {
	GetID() int
	GetPoints() float64
	SetPoints(points float64)
}

// Score is a default struct used as an example of how structs can be implemented for tournify
type Score struct {
	ID     int
	Points float64 // We want to support any type of game where points can be very high or even just decimals
}

// GetID returns the id of the score
func (s *Score) GetID() int {
	return s.ID
}

// GetPoints returns the point value of the score
func (s *Score) GetPoints() float64 {
	return s.Points
}

// SetPoints allows you to set points for the score
func (s *Score) SetPoints(points float64) {
	s.Points = points
}
