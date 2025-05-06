// Package auth provides authentication-related utilities including password hashing and verification.
package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword takes a plain-text password and returns its bcrypt hash.
//
// The hashing algorithm uses bcrypt with the default cost (work factor).
// This hash should be stored securely and never exposed.
//
// Returns:
//   - The bcrypt hash as a string
//   - An error if hashing fails
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a plain-text password with its stored bcrypt hash.
//
// Returns true if the password matches the hash, otherwise returns false.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
