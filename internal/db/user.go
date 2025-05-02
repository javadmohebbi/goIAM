package db

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username      string `gorm:"uniqueIndex;not null"`
	Email         string `gorm:"uniqueIndex"`
	EmailVerified bool   `gorm:"default:false"`
	PhoneNumber   string `gorm:"uniqueIndex"`
	PhoneVerified bool   `gorm:"default:false"`
	PasswordHash  string `gorm:"not null"`

	FirstName  string
	MiddleName string
	LastName   string
	Address    string

	IsActive bool `gorm:"default:true"`

	// Relationships
	Groups   []Group  `gorm:"many2many:user_groups;"`
	Roles    []Role   `gorm:"many2many:user_roles;"`
	Policies []Policy `gorm:"many2many:user_policies;"`

	TOTPSecret    string `gorm:""`
	Requires2FA   bool   `gorm:"default:false"`
	TwoFAVerified bool   `gorm:"-:all"` // runtime only
	BackupCodes   []BackupCode
}

type BackupCode struct {
	gorm.Model
	UserID   uint
	CodeHash string `gorm:"not null"`
	Used     bool   `gorm:"default:false"`
}
