package models

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	Name  string
	label string
}

type PermissionRole struct {
	gorm.Model
	Permission   Permission
	Role         Role
	PermissionID uint
	RoleID       uint
}
