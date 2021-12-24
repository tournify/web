package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title       string
	Slug        string
	Content     string // LongText
	Keywords    string
	Description string
	CreatorID   uint
	Creator     User
}
