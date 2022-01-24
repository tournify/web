package models

import (
	"fmt"
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	Name string `json:"name"`
	Slug string `gorm:"uniqueIndex;size:256;" json:"slug"`
	// TODO Eliminated is specific to one tournament, in the future we could move this to a separate table so that teams are not tournament specific
	Eliminated  int    `json:"-"`
	Keywords    string `json:"-"`
	Description string `json:"-"`
	Games       []Game `gorm:"many2many:game_teams;" json:"-"`
}

// GetID returns the id of the score
func (t *Team) GetID() int {
	return int(t.ID)
}

// Print writes team details to stdout
func (t *Team) Print() string {
	return fmt.Sprintf("Team ID: %d\n", t.GetID())
}

// GetEliminatedCount gets the number of times the team has been eliminated in a multiple elimination tournament
func (t *Team) GetEliminatedCount() int {
	return t.Eliminated
}

func (t *Team) SetEliminatedCount(c int) {
	t.Eliminated = c
}
