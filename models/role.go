package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name  string `gorm:"uniqueIndex;size:256;"`
	Label string `gorm:"uniqueIndex;size:256;"`
	Users []User
}
