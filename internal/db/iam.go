package db

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name     string   `gorm:"uniqueIndex;not null"`
	Users    []User   `gorm:"many2many:user_groups;"`
	Policies []Policy `gorm:"many2many:group_policies;"`
}

type Role struct {
	gorm.Model
	Name     string   `gorm:"uniqueIndex;not null"`
	Users    []User   `gorm:"many2many:user_roles;"`
	Policies []Policy `gorm:"many2many:role_policies;"`
}

type Policy struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;not null"`
	Description string
}
