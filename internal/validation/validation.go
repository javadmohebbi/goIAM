package validation

import (
	"regexp"
	"strings"
)

// ValidateEmail validates the email format based on the configured regex pattern.
// Falls back to checking for the presence of "@" if no regex is set.
func (v *Validation) ValidateEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	if v.cfg.Validation.EmailRegex == "" {
		return strings.Contains(email, "@")
	}
	return regexp.MustCompile(v.cfg.Validation.EmailRegex).MatchString(email)
}

// ValidatePhone validates a phone number using the configured regex pattern.
// Falls back to a minimum length check if no regex is provided.
func (v *Validation) ValidatePhone(phone string) bool {
	if v.cfg.Validation.PhoneRegex == "" {
		return len(phone) >= 6
	}
	return regexp.MustCompile(v.cfg.Validation.PhoneRegex).MatchString(phone)
}

// ValidatePassword checks if the password meets the minimum length
// and matches the configured regex pattern (if set).
func (v *Validation) ValidatePassword(password string) bool {
	if len(password) < v.cfg.Validation.PasswordMinLength {
		return false
	}
	if v.cfg.Validation.PasswordRegex != "" {
		return regexp.MustCompile(v.cfg.Validation.PasswordRegex).MatchString(password)
	}
	return true
}

// ValidateLength returns true if the string's length is within the given range [min, max].
func ValidateLength(value string, min, max int) bool {
	l := len(value)
	return l >= min && l <= max
}

// ValidateWebsite checks whether a URL matches the configured website regex pattern.
// Falls back to checking for "http://" or "https://" prefix if no pattern is set.
func (v *Validation) ValidateWebsite(url string) bool {
	if v.cfg.Validation.WebsiteRegex == "" {
		return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
	}
	return regexp.MustCompile(v.cfg.Validation.WebsiteRegex).MatchString(url)
}
