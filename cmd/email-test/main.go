package main

import (
	"fmt"
	"log"
	"time"

	"github.com/javadmohebbi/goIAM/internal/config"
	"github.com/javadmohebbi/goIAM/internal/db"
	"github.com/javadmohebbi/goIAM/internal/smtpclient"
)

func main() {
	// Load application config
	cfg, err := config.LoadConfig("/Users/javad/Projects/goIAM/config.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	// Dummy user for testing
	user := db.User{
		Username:  "javad",
		FirstName: "Javad",
		Email:     "javad@example.com",
	}

	if err := sendTestEmail(user, cfg); err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	log.Println("Test email sent successfully.")
}

func sendTestEmail(user db.User, cfg *config.Config) error {
	placeholders := map[string]string{
		"Name":    user.FirstName,
		"AppName": cfg.AppName,
		"Year":    fmt.Sprintf("%d", time.Now().Year()),
		"Token":   "test-token-123",
	}

	return smtpclient.SendEmailFromHTMLTemplate(cfg, "Activate Your Account",
		[]string{user.Email}, "/Users/javad/Projects/goIAM/html-templates/reset-password.html", placeholders)
}
