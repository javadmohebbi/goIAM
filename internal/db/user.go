package db

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string `gorm:"uniqueIndex;not null"`
	Email        string `gorm:"uniqueIndex"`
	PasswordHash string `gorm:"not null"`
	IsActive     bool   `gorm:"default:true"`
}
