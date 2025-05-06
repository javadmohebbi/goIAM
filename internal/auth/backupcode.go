// Package auth provides functions for generating and verifying backup codes
// used in two-factor authentication (2FA) or account recovery workflows.
package auth

import (
	"crypto/rand"
	"encoding/base32"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// GenerateBackupCodes creates `n` secure backup codes and their bcrypt hashes.
//
// Each code is a randomly generated 5-byte value encoded using base32 (lowercase, no padding).
// The hashed version of each code is also returned, suitable for secure storage.
//
// Returns:
//   - A slice of plain text backup codes for the user
//   - A slice of bcrypt-hashed codes for server-side storage
//   - An error if random generation or hashing fails
func GenerateBackupCodes(n int) ([]string, []string, error) {
	codes := []string{}  // Plain backup codes to show user
	hashed := []string{} // Corresponding bcrypt-hashed codes for storage

	for i := 0; i < n; i++ {
		buf := make([]byte, 5) // 5 random bytes for each code
		if _, err := rand.Read(buf); err != nil {
			return nil, nil, err
		}

		// Encode to base32, make lowercase, and strip padding
		code := strings.ToLower(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(buf))
		codes = append(codes, code)

		// Generate bcrypt hash of the code
		hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
		if err != nil {
			return nil, nil, err
		}
		hashed = append(hashed, string(hash))
	}
	return codes, hashed, nil
}

// CheckBackupCode verifies whether the provided `code` matches the given bcrypt `hash`.
//
// Returns true if the code is valid (i.e., the hash matches), false otherwise.
func CheckBackupCode(code string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(code)) == nil
}
