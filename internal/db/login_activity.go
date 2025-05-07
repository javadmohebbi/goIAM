// Package db defines database models and operations used by the goIAM service.
package db

import "gorm.io/gorm"

// LoginActivity represents an audit log entry for a user's login event.
// This model is stored in a separate audit log database to record metadata
// such as IP address, browser, operating system, and device information.
type LoginActivity struct {
	gorm.Model

	// UserID is the foreign key reference to the user who logged in.
	UserID uint

	// Username is the user's unique identifier at the time of login.
	Username string

	// IP is the IP address from which the login was performed.
	IP string

	// UserAgent is the full User-Agent string from the HTTP request header.
	UserAgent string

	// OS is the operating system extracted from the User-Agent string.
	OS string

	// Browser is the browser name extracted from the User-Agent string.
	Browser string

	// Device represents the platform or hardware type, e.g., mobile, desktop.
	Device string

	// Status describes the result of the login attempt (e.g., "success", "invalid_password").
	Status string

	// Success is true if the login attempt succeeded; false otherwise.
	Success bool

	// Location is the optional geographical location of the IP address.
	Location string
}
