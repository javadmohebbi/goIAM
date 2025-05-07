package db

import (
	"fmt"

	"gorm.io/gorm"
)

// UpdateProfile updates the user's non-sensitive profile information.
//
// Note: This method does not allow changes to OrganizationID, PasswordHash, or 2FA settings.
// Those should be handled by dedicated methods.
func (u *User) UpdateProfile(db *gorm.DB, updates map[string]interface{}) error {
	// Prevent modification of restricted fields
	delete(updates, "organization_id")
	delete(updates, "password_hash")

	return db.Model(u).Updates(updates).Error
}

// UpdatePasswordHash updates the user's password hash.
//
// This method should be used after securely hashing the password.
// It only updates the PasswordHash field.
func (u *User) UpdatePasswordHash(db *gorm.DB, newHash string) error {
	u.PasswordHash = newHash
	return db.Model(u).Update("password_hash", newHash).Error
}

// Enable2FA sets the user's 2FA secret and marks 2FA as required.
//
// This should be called after successfully verifying the TOTP secret.
func (u *User) Enable2FA(db *gorm.DB, hashedSecret string) error {
	u.Requires2FA = true
	u.TOTPSecret = hashedSecret
	return db.Model(u).Updates(map[string]interface{}{
		"requires_2fa": true,
		"totp_secret":  hashedSecret,
	}).Error
}

// Disable2FA disables two-factor authentication and clears the secret.
//
// This should be called when a user intentionally disables 2FA.
func (u *User) Disable2FA(db *gorm.DB) error {
	u.Requires2FA = false
	u.TOTPSecret = ""
	return db.Model(u).Updates(map[string]interface{}{
		"requires_2fa": false,
		"totp_secret":  "",
	}).Error
}

// UserAccessSummary holds the list of roles, groups, and policy names associated with a user.
type UserAccessSummary struct {
	UserID   uint
	Username string
	Roles    []string
	Groups   []string
	Policies []string
}

// GetUserAccessSummary retrieves all roles, groups, and policies assigned to the user.
func (u *User) GetUserAccessSummary(db *gorm.DB) (*UserAccessSummary, error) {
	if err := db.Preload("Roles").Preload("Groups").Preload("Policies").First(u, u.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load user relationships: %w", err)
	}

	summary := &UserAccessSummary{
		UserID:   u.ID,
		Username: u.Username,
	}

	for _, role := range u.Roles {
		summary.Roles = append(summary.Roles, role.Name)
	}
	for _, group := range u.Groups {
		summary.Groups = append(summary.Groups, group.Name)
	}
	for _, policy := range u.Policies {
		summary.Policies = append(summary.Policies, policy.Name)
	}

	return summary, nil
}
