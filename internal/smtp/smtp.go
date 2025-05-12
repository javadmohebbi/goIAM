package smtpclient

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strings"

	"github.com/javadmohebbi/goIAM/internal/config"
	// adjust this import to your actual config package
)

// SendEmail sends an email with specified MIME type (e.g., text/plain or text/html)
func SendEmail(cfg *config.Config, subject, to, body, mime string) error {
	auth := smtp.PlainAuth("", cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Host)

	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", cfg.SMTP.FromName, cfg.SMTP.FromEmail)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = mime + "; charset=\"UTF-8\""

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n" + body)

	addr := fmt.Sprintf("%s:%d", cfg.SMTP.Host, cfg.SMTP.Port)
	return smtp.SendMail(addr, auth, cfg.SMTP.FromEmail, []string{to}, []byte(msg.String()))
}

// SendEmailFromHTMLTemplate reads an HTML file and replaces {{PLACEHOLDER}}s with provided values
func SendEmailFromHTMLTemplate(cfg *config.Config, subject, to, templatePath string, placeholders map[string]string) error {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return err
	}

	tmpl, err := template.New("email").Parse(string(content))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, placeholders)
	if err != nil {
		return err
	}

	return SendEmail(cfg, subject, to, buf.String(), "text/html")
}
