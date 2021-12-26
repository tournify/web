package models

import (
	"gorm.io/gorm"
)

const TournamentPrivacyPublic = 1
const TournamentPrivacyPrivate = 2
const TournamentPrivacyLink = 3

type Tournament struct {
	gorm.Model
	Name        string
	Slug        string `gorm:"uniqueIndex;size:256;"`
	Keywords    string
	Description string
	Type        int
	Privacy     int // Public, Private, Anyone with link
	Options     []TournamentOption
	Users       []TournamentUser
}

type TournamentOption struct {
	gorm.Model
	Key          string
	Value        string
	TournamentID uint
	Tournament   Tournament
}

type TournamentUser struct {
	Tournament   Tournament
	User         User
	TournamentID uint
	UserID       uint
}
