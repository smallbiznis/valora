package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/smtp"
	"path/filepath"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type SMTPProvider struct {
	cfg Config
}

func NewSMTP(cfg Config) *SMTPProvider {
	return &SMTPProvider{cfg: cfg}
}

func (p *SMTPProvider) Send(ctx context.Context, to []string, subject string, htmlBody string) error {
	auth := smtp.PlainAuth("", p.cfg.Username, p.cfg.Password, p.cfg.Host)
	addr := fmt.Sprintf("%s:%d", p.cfg.Host, p.cfg.Port)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n%s\r\n%s", to[0], subject, mime, htmlBody))

	return smtp.SendMail(addr, auth, p.cfg.From, to, msg)
}

func (p *SMTPProvider) SendTemplate(ctx context.Context, to []string, templateName string, data interface{}) error {
	// Locate template file (Assumes running from project root or configured path)
	// For MVP, using relative path.
	tmplPath := filepath.Join("internal", "providers", "email", "templates", templateName+".html")

	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Extract subject from data if provided, otherwise use template-based default
	subject := "Notification from Valora"
	if dataMap, ok := data.(map[string]interface{}); ok {
		if subj, exists := dataMap["subject"]; exists {
			if subjStr, ok := subj.(string); ok {
				subject = subjStr
			}
		} else {
			// Use template-specific defaults
			switch templateName {
			case "invite_member":
				if orgName, ok := dataMap["org_name"].(string); ok {
					subject = fmt.Sprintf("You're invited to join %s", orgName)
				} else {
					subject = "You're invited to join a team"
				}
			case "invoice_new":
				subject = "New invoice from Valora"
			}
		}
	}

	return p.Send(ctx, to, subject, body.String())
}
