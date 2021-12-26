package models

import "gorm.io/gorm"

type Score struct {
	gorm.Model
	Game   Game    `json:"-"`
	GameID uint    `json:"-"`
	Team   Team    `json:"-"`
	TeamID uint    `json:"-"`
	Score  float64 `json:"score"`
}

// GetID returns the id of the score
func (s *Score) GetID() int {
	return int(s.ID)
}

// GetPoints returns the point value of the score
func (s *Score) GetPoints() float64 {
	return s.Score
}

// SetPoints allows you to set points for the score
func (s *Score) SetPoints(points float64) {
	s.Score = points
}
