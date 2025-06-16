package config

import "context"

type contextKey string

const configKey contextKey = "config"

func WithConfig(ctx context.Context, cfg Config) context.Context {
	return context.WithValue(ctx, configKey, cfg)
}

func FromContext(ctx context.Context) *Config {
	if cfg, ok := ctx.Value(configKey).(Config); ok {
		return &cfg
	}
	if cfg, ok := ctx.Value(configKey).(*Config); ok {
		return cfg
	}
	return nil
}
