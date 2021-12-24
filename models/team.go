package models

import "gorm.io/gorm"

type Team struct {
	gorm.Model
	Name        string
	Slug        string
	Keywords    string
	Description string
}
