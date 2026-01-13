package email

import (
	"os"
	"strconv"

	"go.uber.org/fx"
)

var Module = fx.Module("providers.email",
	fx.Provide(NewFromEnv),
)

func NewFromEnv() Provider {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	cfg := Config{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     port,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
	}
	// Defaults for dev
	// Defaults for dev (MailHog)
	// If Host is empty OR (Host is localhost and Port is 0/unset) -> Use MailHog defaults
	if cfg.Host == "" || (cfg.Host == "localhost" && cfg.Port == 0) {
		cfg.Host = "localhost"
		cfg.Port = 1025 // Mailhog default
		cfg.From = "no-reply@valora.test"
	}
	return NewSMTP(cfg)
}
