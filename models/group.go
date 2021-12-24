package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name         string
	Slug         string
	Tournament   Tournament
	TournamentID uint
}
