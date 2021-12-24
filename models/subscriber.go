package models

import "gorm.io/gorm"

// 1	id	int unsigned	NULL	NULL	NO	NULL	auto_increment
//2	email	varchar(255)	utf8	utf8_unicode_ci	NO	NULL
//3	created_at	timestamp	NULL	NULL	NO	0000-00-00 00:00:00
//4	updated_at	timestamp	NULL	NULL	NO	0000-00-00 00:00:00
//5	deleted_at	timestamp	NULL	NULL	YES	NULL

type Subscriber struct {
	gorm.Model
	Email string
}
