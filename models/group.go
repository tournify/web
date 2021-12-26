package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name         string
	Slug         string `gorm:"uniqueIndex;size:256;"`
	Tournament   Tournament
	TournamentID uint
}
