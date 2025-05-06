// Package db defines user-related models for authentication, identity management,
// and 2FA (Two-Factor Authentication) features in the IAM system.
package db

import "gorm.io/gorm"

// User represents an account in the system with identity and access management attributes.
//
// Fields:
//   - Username, Email, and PhoneNumber uniquely identify the user
//   - EmailVerified and PhoneVerified indicate verification status
//   - PasswordHash stores the user's hashed password
//   - FirstName, MiddleName, LastName, and Address hold personal info
//   - IsActive controls if the user account is currently enabled
//   - Groups, Roles, and Policies are used for access control (many-to-many)
//   - TOTPSecret and BackupCodes support 2FA functionality
type User struct {
	gorm.Model
	Username      string `gorm:"uniqueIndex;not null"` // Unique username
	Email         string `gorm:"uniqueIndex"`          // Optional unique email
	EmailVerified bool   `gorm:"default:false"`        // Email verification flag
	PhoneNumber   string // Optional phone number
	PhoneVerified bool   `gorm:"default:false"` // Phone number verification flag
	PasswordHash  string `gorm:"not null"`      // Bcrypt-hashed password

	FirstName  string // User's first name
	MiddleName string // User's middle name
	LastName   string // User's last name
	Address    string // Mailing or home address

	IsActive bool `gorm:"default:true"` // Whether the account is enabled

	// Relationships
	Groups   []Group  `gorm:"many2many:user_groups;"`   // Group memberships
	Roles    []Role   `gorm:"many2many:user_roles;"`    // Assigned roles
	Policies []Policy `gorm:"many2many:user_policies;"` // Directly attached policies

	TOTPSecret    string       // TOTP secret used for 2FA
	Requires2FA   bool         `gorm:"default:false"` // Whether 2FA is required
	TwoFAVerified bool         `gorm:"-:all"`         // Set at runtime only (ignored by GORM)
	BackupCodes   []BackupCode // List of backup codes for 2FA recovery
}

// BackupCode stores a one-time-use code for users who enable 2FA.
//
// Fields:
//   - UserID: foreign key reference to the User
//   - CodeHash: bcrypt-hashed backup code
//   - Used: whether the code has already been consumed
type BackupCode struct {
	gorm.Model
	UserID   uint   // Foreign key to User
	CodeHash string `gorm:"not null"`      // Hashed version of the backup code
	Used     bool   `gorm:"default:false"` // Whether the code has been used
}
