package models

import "gorm.io/gorm"

// 1	id	int unsigned	NULL	NULL	NO	NULL	auto_increment
//2	name	varchar(255)	utf8	utf8_unicode_ci	NO	NULL
//3	slug	varchar(255)	utf8	utf8_unicode_ci	NO	NULL
//4	keywords	varchar(255)	utf8	utf8_unicode_ci	NO	NULL
//5	description	varchar(255)	utf8	utf8_unicode_ci	NO	NULL
//6	created_at	timestamp	NULL	NULL	NO	0000-00-00 00:00:00
//7	updated_at	timestamp	NULL	NULL	NO	0000-00-00 00:00:00
//8	deleted_at	timestamp	NULL	NULL	YES	NULL

type Player struct {
	gorm.Model
	Name        string
	Slug        string `gorm:"uniqueIndex;size:256;"`
	Keywords    string
	Description string
}
