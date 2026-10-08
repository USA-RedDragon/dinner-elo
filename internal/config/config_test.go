package config_test

import (
	"errors"
	"testing"

	"github.com/USA-RedDragon/configulator/v2"
	"github.com/USA-RedDragon/dinner-elo/internal/config"
)

func validConfig(t *testing.T) config.Config {
	t.Helper()
	cfg, err := configulator.New(config.ConfigSchema()).Default()
	if err != nil {
		t.Fatalf("failed to build default config: %v", err)
	}
	cfg.HTTP.URL = "http://localhost:8080"
	cfg.Auth.JWTSecret = "secret"
	return cfg
}

func TestDefaults(t *testing.T) {
	t.Parallel()
	cfg, err := configulator.New(config.ConfigSchema()).Default()
	if err != nil {
		t.Fatalf("failed to build default config: %v", err)
	}
	if cfg.LogLevel != config.LogLevelInfo {
		t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, config.LogLevelInfo)
	}
	if cfg.Storage.Type != config.StorageTypeSQLite {
		t.Errorf("Storage.Type = %q, want %q", cfg.Storage.Type, config.StorageTypeSQLite)
	}
	if cfg.HTTP.Port != 8080 {
		t.Errorf("HTTP.Port = %d, want 8080", cfg.HTTP.Port)
	}
	if cfg.Metrics.Port != 9000 {
		t.Errorf("Metrics.Port = %d, want 9000", cfg.Metrics.Port)
	}
	if cfg.PProf.Port != 9999 {
		t.Errorf("PProf.Port = %d, want 9999", cfg.PProf.Port)
	}
}

func TestValidate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		modify func(*config.Config)
		want   error
	}{
		{name: "valid", modify: func(*config.Config) {}},
		{name: "bad log level", modify: func(c *config.Config) { c.LogLevel = "trace" }, want: config.ErrInvalidLogLevel},
		{name: "bad http port", modify: func(c *config.Config) { c.HTTP.Port = 0 }, want: config.ErrInvalidHTTPPort},
		{name: "bad http address", modify: func(c *config.Config) { c.HTTP.Address = "nope" }, want: config.ErrInvalidHTTPAddress},
		{name: "bad trusted proxy", modify: func(c *config.Config) { c.HTTP.TrustedProxies = []string{"nope"} }, want: config.ErrInvalidTrustedProxies},
		{name: "trusted proxy ip", modify: func(c *config.Config) { c.HTTP.TrustedProxies = []string{"10.0.0.1"} }},
		{name: "bad metrics port", modify: func(c *config.Config) { c.Metrics.Enabled = true; c.Metrics.Port = 70000 }, want: config.ErrInvalidMetricsPort},
		{name: "metrics without address", modify: func(c *config.Config) { c.Metrics.Enabled = true }, want: config.ErrInvalidMetricsAddress},
		{name: "bad pprof port", modify: func(c *config.Config) { c.PProf.Enabled = true; c.PProf.Port = -1 }, want: config.ErrInvalidPProfPort},
		{name: "pprof without address", modify: func(c *config.Config) { c.PProf.Enabled = true }, want: config.ErrInvalidPProfAddress},
		{name: "bad storage type", modify: func(c *config.Config) { c.Storage.Type = "redis" }, want: config.ErrInvalidStorageType},
		{name: "empty dsn", modify: func(c *config.Config) { c.Storage.DSN = "" }, want: config.ErrInvalidStorageDSN},
		{name: "empty jwt secret", modify: func(c *config.Config) { c.Auth.JWTSecret = "" }, want: config.ErrInvalidAuthJWTSecret},
		{name: "empty http url", modify: func(c *config.Config) { c.HTTP.URL = "" }, want: config.ErrInvalidHTTPURL},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := validConfig(t)
			tt.modify(&cfg)
			if err := cfg.Validate(); !errors.Is(err, tt.want) {
				t.Errorf("Validate() = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestValidateTrustedProxies(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		proxies []string
		want    error
	}{
		{name: "ipv4", proxies: []string{"10.0.0.1"}},
		{name: "ipv6", proxies: []string{"::1"}},
		{name: "cidr v4", proxies: []string{"10.0.0.0/8"}},
		{name: "cidr v6", proxies: []string{"fd00::/8"}},
		{name: "mixed", proxies: []string{"172.16.0.0/12", "192.168.1.1", "2001:db8::/32", "fe80::1"}},
		{name: "invalid", proxies: []string{"not-an-ip"}, want: config.ErrInvalidTrustedProxies},
		{name: "invalid cidr prefix", proxies: []string{"10.0.0.0/33"}, want: config.ErrInvalidTrustedProxies},
		{name: "invalid after valid", proxies: []string{"192.168.0.0/16", "300.1.1.1"}, want: config.ErrInvalidTrustedProxies},
		{name: "empty entry", proxies: []string{""}, want: config.ErrInvalidTrustedProxies},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := validConfig(t)
			cfg.HTTP.TrustedProxies = tt.proxies
			if err := cfg.Validate(); !errors.Is(err, tt.want) {
				t.Errorf("Validate() = %v, want %v", err, tt.want)
			}
		})
	}
}
