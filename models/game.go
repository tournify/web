package models

import "gorm.io/gorm"

type Game struct {
	gorm.Model
	Name         string
	Slug         string
	Keywords     string
	Description  string
	Type         int
	Tournament   Tournament
	TournamentID uint
}
