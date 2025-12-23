package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds application configuration.
type Config struct {
	AppName      string
	AppVersion   string
	Mode         string
	Environment  string

	OTLPEndpoint string

	Cloud CloudConfig

	DBType            string
	DBHost            string
	DBPort            string
	DBName            string
	DBUser            string
	DBPassword        string
	DBSSLMode         string
	DBMaxIdleConn     int
	DBMaxOpenConn     int
	DBConnMaxLifetime int
	DBConnMaxIdleTime int
}

type CloudConfig struct {
	OrganizationID string
	OrganizationName string
	Metrics        CloudMetricsConfig
}

type CloudMetricsConfig struct {
	Enabled   bool
	Exporter  string
	Endpoint  string
	AuthToken string
}

// Load loads configuration from environment variables and .env file.
func Load() Config {
	_ = godotenv.Load()

	mode := normalizeMode(getenv("APP_MODE", ModeOSS))
	cfg := Config{
		AppName:      getenv("APP_SERVICE", "valora"),
		AppVersion:   getenv("APP_VERSION", "0.1.0"),
		Mode:         mode,
		Environment:  getenv("ENVIRONMENT", "development"),
		OTLPEndpoint: getenv("OTLP_ENDPOINT", "localhost:4317"),
		Cloud: CloudConfig{
			OrganizationID: strings.TrimSpace(getenv("CLOUD_ORGANIZATION_ID", "")),
			OrganizationName: getenv("CLOUD_ORGANIZATION_NAME", ""),
			Metrics: CloudMetricsConfig{
				Enabled:   getenvBool("CLOUD_METRICS_ENABLED", true),
				Exporter:  strings.ToLower(getenv("CLOUD_METRICS_EXPORTER", "")),
				Endpoint:  strings.TrimSpace(getenv("CLOUD_METRICS_ENDPOINT", "")),
				AuthToken: strings.TrimSpace(getenv("CLOUD_METRICS_AUTH_TOKEN", "")),
			},
		},
		DBType:       getenv("DATABASE_TYPE", "postgres"),
		DBHost:       getenv("DATABASE_HOST", "localhost"),
		DBPort:       getenv("DATABASE_PORT", "5433"),
		DBName:       getenv("DATABASE_NAME", "postgres"),
		DBUser:       getenv("DATABASE_USER", "postgres"),
		DBPassword:   getenv("DATABASE_PASSWORD", "35411231"),
		DBSSLMode:    getenv("DATABASE_SSLMODE", "disable"),
	}

	return cfg
}

const (
	ModeOSS   = "oss"
	ModeCloud = "cloud"
	ModeStandalone = "standalone"
)

func (c Config) IsCloud() bool {
	return c.Mode == ModeCloud
}

func normalizeMode(raw string) string {
	value := strings.ToLower(strings.TrimSpace(raw))
	switch value {
	case ModeCloud:
		return ModeCloud
	case ModeStandalone, ModeOSS:
		return ModeOSS
	default:
		return ModeOSS
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvBool(key string, def bool) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	if value == "" {
		return def
	}
	switch value {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return def
	}
}

func parseServices(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		out = append(out, p)
	}
	if len(out) == 0 {
		log.Println("no services enabled for migration")
	}
	return out
}
