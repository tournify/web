package models

import "gorm.io/gorm"

type Score struct {
	gorm.Model
	Game   Game
	GameID uint
	Team   Team
	TeamID uint
	Score  float64
}
