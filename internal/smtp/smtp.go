package smtpclient

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strings"

	"github.com/javadmohebbi/goIAM/internal/config"
)

// SendEmail sends an email using the provided SMTP configuration.
// You can specify the subject, recipient list, body content, and MIME type.
// MIME should be either "text/plain" or "text/html" for most use cases.
func SendEmail(cfg *config.Config, subject string, to []string, body, mime string) error {
	auth := smtp.PlainAuth("", cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Host)

	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", cfg.SMTP.FromName, cfg.SMTP.FromEmail)
	headers["To"] = strings.Join(to, ", ")
	headers["Subject"] = subject
	headers["Reply-To"] = cfg.SMTP.FromEmail
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = mime + "; charset=\"UTF-8\""
	headers["Content-Transfer-Encoding"] = "quoted-printable"

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n" + body)

	addr := fmt.Sprintf("%s:%d", cfg.SMTP.Host, cfg.SMTP.Port)
	return smtp.SendMail(addr, auth, cfg.SMTP.FromEmail, to, []byte(msg.String()))
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
		return fmt.Errorf("unable to read template: %w", err)
	}

	tmpl, err := template.New("email").Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, placeholders)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return SendEmail(cfg, subject, to, buf.String(), "text/html")
}
