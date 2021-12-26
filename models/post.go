package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Localizations []PostLocalization
	Slug          string `gorm:"uniqueIndex;size:256;"`
	CreatorID     uint
	Creator       User
}

type PostLocalization struct {
	gorm.Model
	PostID      uint
	Post        Post
	Lang        string
	Content     string // LongText
	Keywords    string
	Description string
	Title       string
}
