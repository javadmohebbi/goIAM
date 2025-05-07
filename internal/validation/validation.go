package validation

import (
	"regexp"
	"strings"

	"github.com/javadmohebbi/goIAM/internal/config"
)

// ValidateEmail checks if the email matches the regex pattern defined in the config.
func ValidateEmail(cfg *config.Config, email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	if cfg.Validation.EmailRegex == "" {
		return strings.Contains(email, "@")
	}
	return regexp.MustCompile(cfg.Validation.EmailRegex).MatchString(email)
}

// ValidatePhone checks if the phone number matches the regex pattern defined in the config.
func ValidatePhone(cfg *config.Config, phone string) bool {
	if cfg.Validation.PhoneRegex == "" {
		return len(phone) >= 6
	}
	return regexp.MustCompile(cfg.Validation.PhoneRegex).MatchString(phone)
}

// ValidatePassword checks if the password meets length and complexity regex in config.
func ValidatePassword(cfg *config.Config, password string) bool {
	if len(password) < cfg.Validation.PasswordMinLength {
		return false
	}
	if cfg.Validation.PasswordRegex != "" {
		return regexp.MustCompile(cfg.Validation.PasswordRegex).MatchString(password)
	}
	return true
}

// ValidateLength checks if a value's length is between min and max.
func ValidateLength(value string, min, max int) bool {
	l := len(value)
	return l >= min && l <= max
}

// ValidateWebsite checks if a URL matches the configured website regex pattern.
func ValidateWebsite(cfg *config.Config, url string) bool {
	if cfg.Validation.WebsiteRegex == "" {
		return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
	}
	return regexp.MustCompile(cfg.Validation.WebsiteRegex).MatchString(url)
}
