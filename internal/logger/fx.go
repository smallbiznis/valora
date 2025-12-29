package logger

import (
	"context"
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Config defines logger configuration.
type Config struct {
	Level string
}

// ProvideConfig loads logger config from environment variables.
// Supported env vars:
// - LOG_LEVEL: debug|info|warn|error (default: info)
func ProvideConfig() Config {
	return Config{Level: os.Getenv("LOG_LEVEL")}
}

// NewFromConfig creates a zap logger from Config and replaces globals.
func NewFromConfig(cfg Config) (*zap.Logger, error) {
	return New(cfg.Level)
}

func registerHooks(lc fx.Lifecycle, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			_ = ctx
			_ = log.Sync()
			return nil
		},
	})
}

// Module wires the global zap logger for the application.
var Module = fx.Module("logger",
	fx.Provide(
		ProvideConfig,
		NewFromConfig,
	),
	fx.Invoke(registerHooks),
)
