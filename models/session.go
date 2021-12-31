package models

import (
	"gorm.io/gorm"
	"time"
)

type Session struct {
	gorm.Model
	Identifier  string
	UserID      *uint
	User        *User
	ExpiresAt   time.Time
	Tournaments []Tournament `gorm:"many2many:session_tournaments;" json:"tournaments"`
}

func (s Session) HasExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}
