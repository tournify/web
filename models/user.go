// Package models defines all the database models for the application
package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name        string
	Email       string
	Password    string
	ActivatedAt *time.Time
	Tokens      []Token `gorm:"polymorphic:Model;"`
	Sessions    []Session
	UserGroup   int // Admin, free user, paid basic user, paid whitelabel user
	Tournaments []TournamentUser
}
