package auth

import (
	"crypto/rand"
	"encoding/base32"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func GenerateBackupCodes(n int) ([]string, []string, error) {
	codes := []string{}
	hashed := []string{}

	for i := 0; i < n; i++ {
		buf := make([]byte, 5)
		if _, err := rand.Read(buf); err != nil {
			return nil, nil, err
		}
		code := strings.ToLower(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(buf))
		codes = append(codes, code)

		hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
		if err != nil {
			return nil, nil, err
		}
		hashed = append(hashed, string(hash))
	}
	return codes, hashed, nil
}

func CheckBackupCode(code string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(code)) == nil
}
