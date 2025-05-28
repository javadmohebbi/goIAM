package smtpclient

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/smtp"
	"os"
	"strings"

	"github.com/javadmohebbi/goIAM/internal/config"
)

// SendEmail sends a plain-text email without TLS or authentication.
// Intended for use with a fake/local SMTP server in testing environments.
func SendEmail(cfg *config.Config, subject string, to []string, body, mime string) error {
	headers := map[string]string{
		"From":                      fmt.Sprintf("%s <%s>", cfg.SMTP.FromName, cfg.SMTP.FromEmail),
		"To":                        strings.Join(to, ", "),
		"Subject":                   subject,
		"Reply-To":                  cfg.SMTP.FromEmail,
		"MIME-Version":              "1.0",
		"Content-Type":              mime + "; charset=\"UTF-8\"",
		"Content-Transfer-Encoding": "quoted-printable",
	}

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n" + body)

	// both for ipv4 and ipv6
	addr := net.JoinHostPort(cfg.SMTP.Host, fmt.Sprintf("%d", cfg.SMTP.Port))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}

	client, err := smtp.NewClient(conn, cfg.SMTP.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	if err := client.Hello("localhost"); err != nil {
		return fmt.Errorf("HELO failed: %w", err)
	}

	if err := client.Mail(cfg.SMTP.FromEmail); err != nil {
		return fmt.Errorf("MAIL FROM failed: %w", err)
	}
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("RCPT TO failed for %s: %w", recipient, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA command failed: %w", err)
	}
	if _, err := w.Write([]byte(msg.String())); err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to complete DATA: %w", err)
	}

	if err := client.Quit(); err != nil {
		return fmt.Errorf("failed to close SMTP session: %w", err)
	}

	return nil
}

// SendPlainTextEmail is a helper that sends a plain text email using UTF-8 encoding.
func SendPlainTextEmail(cfg *config.Config, subject string, to []string, body string) error {
	return SendEmail(cfg, subject, to, body, "text/plain")
}

// SendEmailFromHTMLTemplate reads an HTML file, replaces {{.Key}} placeholders with the
// provided values, and sends it as an HTML email to the specified recipients.
//
// Parameters:
//   - cfg: SMTP configuration (loaded from your app config)
//   - subject: Email subject
//   - to: List of recipient email addresses
//   - templatePath: Path to the HTML template file
//   - placeholders: Map of key-value pairs to replace placeholders in the HTML template
//
// The HTML file should contain Go template-style placeholders like {{.Name}}, {{.Link}}, etc.
//
// Example:
//
//	err := smtpclient.SendEmailFromHTMLTemplate(cfg, "Welcome",
//		[]string{"user@example.com"}, "templates/welcome.html", map[string]string{
//			"Name": "Javad",
//			"Link": "https://goiam.io/verify",
//		})
//	if err != nil {
//		log.Println("Send failed:", err)
//	}
func SendEmailFromHTMLTemplate(cfg *config.Config, subject string, to []string, templatePath string, placeholders map[string]string) error {

	content, err := os.ReadFile(templatePath)
	if err != nil {
		log.Println("unable to read template:", err)
		return fmt.Errorf("unable to read template: %w", err)
	}

	tmpl, err := template.New("email").Parse(string(content))
	if err != nil {
		log.Println("failed to parse template:", err)
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, placeholders)
	if err != nil {
		log.Println("failed to execute template:", err)
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return SendEmail(cfg, subject, to, buf.String(), "text/html")
}
