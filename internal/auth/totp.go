// Package auth provides authentication utilities, including Time-based One-Time Password (TOTP)
// generation and validation using the pquerna/otp library.
package auth

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// GenerateTOTPSecret generates a new TOTP secret key for a user.
//
// The generated key can be used to configure an authenticator app (like Google Authenticator).
// It also returns a provisioning URI (qrURL) that can be encoded into a QR code for scanning.
//
// Parameters:
//   - username: the account or email to associate with the TOTP key
//   - issuer: the organization or service name (shown in authenticator apps)
//
// Returns:
//   - key: the generated *otp.Key object
//   - qrURL: the provisioning URI (e.g., otpauth://totp/...)
//   - err: any error that occurred during key generation
func GenerateTOTPSecret(username, issuer string) (key *otp.Key, qrURL string, err error) {
	key, err = totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: username,
	})
	if err != nil {
		return nil, "", err
	}
	return key, key.URL(), nil
}

// ValidateTOTP validates a TOTP code against the shared secret.
//
// This function should be used to verify the code submitted by the user during 2FA login.
//
// Parameters:
//   - secret: the shared TOTP secret
//   - code: the 6-digit code from the authenticator app
//
// Returns true if the code is valid for the current time window.
func ValidateTOTP(secret, code string) bool {
	return totp.Validate(code, secret)
}
