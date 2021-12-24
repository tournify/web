package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name  string
	Label string
}

type RoleUser struct {
	gorm.Model
	User   User
	UserID uint
}
