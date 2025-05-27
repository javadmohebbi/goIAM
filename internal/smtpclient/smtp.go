package smtpclient

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/smtp"
	"os"
	"strings"

	"github.com/javadmohebbi/goIAM/internal/config"
)

// SendEmail sends an email using STARTTLS if the server requires encryption.
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

	addr := fmt.Sprintf("%s:%d", cfg.SMTP.Host, cfg.SMTP.Port)
	tlsConfig := &tls.Config{ServerName: cfg.SMTP.Host}
	var client *smtp.Client
	var conn net.Conn
	var err error

	switch cfg.SMTP.Port {
	case 465:
		if !cfg.SMTP.UseTLS {
			return fmt.Errorf("implicit TLS on port 465 requires UseTLS=true")
		}
		fmt.Println("[DEBUG] Connecting with implicit TLS on port 465")
		conn, err = tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to dial TLS: %w", err)
		}
		client, err = smtp.NewClient(conn, cfg.SMTP.Host)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
	case 587:
		fmt.Println("[DEBUG] Connecting over plain TCP on port 587")
		conn, err = net.Dial("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
		client, err = smtp.NewClient(conn, cfg.SMTP.Host)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
		if err = client.Hello("localhost"); err != nil {
			return fmt.Errorf("HELO failed: %w", err)
		}
		if cfg.SMTP.UseTLS {
			if ok, _ := client.Extension("STARTTLS"); ok {
				fmt.Printf("[DEBUG] STARTTLS supported on port %d, proceeding\n", cfg.SMTP.Port)
				if err = client.StartTLS(tlsConfig); err != nil {
					return fmt.Errorf("STARTTLS failed on port %d: %w", cfg.SMTP.Port, err)
				}
				fmt.Printf("[DEBUG] STARTTLS completed on port %d\n", cfg.SMTP.Port)
			} else {
				return fmt.Errorf("STARTTLS not supported on port %d", cfg.SMTP.Port)
			}
		}
	case 25:
		fmt.Println("[DEBUG] Connecting over plain TCP on port 25")
		conn, err = net.Dial("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
		client, err = smtp.NewClient(conn, cfg.SMTP.Host)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
		if err = client.Hello("localhost"); err != nil {
			return fmt.Errorf("HELO failed: %w", err)
		}
		if cfg.SMTP.UseTLS {
			if ok, _ := client.Extension("STARTTLS"); ok {
				fmt.Printf("[DEBUG] STARTTLS supported on port %d, proceeding\n", cfg.SMTP.Port)
				if err = client.StartTLS(tlsConfig); err != nil {
					return fmt.Errorf("STARTTLS failed on port %d: %w", cfg.SMTP.Port, err)
				}
				fmt.Printf("[DEBUG] STARTTLS completed on port %d\n", cfg.SMTP.Port)
			} else {
				return fmt.Errorf("STARTTLS not supported on port %d", cfg.SMTP.Port)
			}
		}
	default:
		return fmt.Errorf("unsupported SMTP port: %d", cfg.SMTP.Port)
	}

	auth := smtp.PlainAuth("", cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Host)
	log.Println(cfg.SMTP.Username)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("auth failed: %w", err)
	}
	fmt.Println("[DEBUG] Authenticated with SMTP server")

	if err = client.Mail(cfg.SMTP.FromEmail); err != nil {
		return err
	}
	fmt.Printf("[DEBUG] MAIL FROM: %s\n", cfg.SMTP.FromEmail)
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return err
		}
		fmt.Printf("[DEBUG] RCPT TO: %s\n", recipient)
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	fmt.Println("[DEBUG] Beginning DATA transmission")
	if _, err = w.Write([]byte(msg.String())); err != nil {
		return err
	}
	fmt.Println("[DEBUG] Email body written")
	if err = w.Close(); err != nil {
		return err
	}
	fmt.Println("[DEBUG] Email transmission completed")

	err = client.Quit()
	if err == nil {
		fmt.Println("[DEBUG] SMTP session ended")
	}
	return err
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
