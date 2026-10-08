package config

//go:generate go tool configulator -type Config

import (
	"errors"
	"fmt"
	"net"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

type Config struct {
	LogLevel LogLevel `name:"log-level" description:"Logging level for the application. One of debug, info, warn, or error" default:"info"`
	Storage  Storage  `name:"storage" description:"Storage configuration"`
	HTTP     HTTP     `name:"http" description:"HTTP server configuration"`
	Metrics  Metrics  `name:"metrics" description:"Metrics server configuration"`
	PProf    PProf    `name:"pprof" description:"PProf server configuration"`
	Auth     Auth     `name:"auth" description:"Authentication configuration"`
}

type Auth struct {
	JWTSecret    string `name:"jwt-secret" secret:"true" description:"JWT secret for signing tokens. Required"`
	ClientID     string `name:"client-id" description:"OAuth2 client ID"`
	ClientSecret string `name:"client-secret" secret:"true" description:"OAuth2 client secret"`
	TokenURL     string `name:"token-url" description:"OAuth2 token URL, e.g. https://example.com/oauth2/token"`
	UserURL      string `name:"user-url" description:"OAuth2 user information URL, e.g. https://example.com/oauth2/userinfo"`
}

type HTTP struct {
	URL            string   `name:"url" description:"URL where the HTTP server is deployed, used for redirects. Required"`
	Address        string   `name:"address" description:"IP address to listen on. Empty listens on all interfaces"`
	Port           int      `name:"port" description:"Port to listen on" default:"8080"`
	TrustedProxies []string `name:"trusted-proxies" description:"Trusted proxies for the HTTP server"`
}

type Metrics struct {
	Enabled bool   `name:"enabled" description:"Enable metrics server"`
	Address string `name:"address" description:"Address to listen on"`
	Port    int    `name:"port" description:"Port to listen on" default:"9000"`
}

type PProf struct {
	Enabled bool   `name:"enabled" description:"Enable pprof server"`
	Address string `name:"address" description:"Address to listen on"`
	Port    int    `name:"port" description:"Port to listen on" default:"9999"`
}

type StorageType string

const (
	StorageTypeMySQL    StorageType = "mysql"
	StorageTypePostgres StorageType = "postgres"
	StorageTypeSQLite   StorageType = "sqlite"
)

type Storage struct {
	Type StorageType `name:"type" description:"Storage type. One of mysql, postgres, sqlite" default:"sqlite"`
	DSN  string      `name:"dsn" description:"Data source name for the storage" default:"file:./.dinner-elo.db?cache=shared&mode=rwc"`
}

var (
	ErrInvalidLogLevel       = errors.New("invalid log level provided")
	ErrInvalidHTTPURL        = errors.New("invalid HTTP URL provided, must not be empty")
	ErrInvalidHTTPPort       = errors.New("invalid HTTP port provided, must be between 1 and 65535")
	ErrInvalidMetricsPort    = errors.New("invalid metrics port provided, must be between 1 and 65535")
	ErrInvalidPProfPort      = errors.New("invalid pprof port provided, must be between 1 and 65535")
	ErrInvalidHTTPAddress    = errors.New("invalid HTTP address provided, must be a valid IP")
	ErrInvalidMetricsAddress = errors.New("invalid metrics address provided, must be a valid IP")
	ErrInvalidPProfAddress   = errors.New("invalid pprof address provided, must be a valid IP")
	ErrInvalidTrustedProxies = errors.New("invalid trusted proxies provided, must be a valid IP or CIDR range")
	ErrInvalidStorageType    = errors.New("invalid storage type provided, must be one of mysql, postgres, sqlite")
	ErrInvalidStorageDSN     = errors.New("invalid storage DSN provided, must not be empty")
	ErrInvalidAuthJWTSecret  = errors.New("invalid JWT secret provided, must not be empty")
)

//nolint:gocyclo
func (c Config) Validate() error {
	if c.LogLevel != LogLevelDebug && c.LogLevel != LogLevelInfo && c.LogLevel != LogLevelWarn && c.LogLevel != LogLevelError {
		return ErrInvalidLogLevel
	}

	if c.HTTP.Port < 1 || c.HTTP.Port > 65535 {
		return ErrInvalidHTTPPort
	}

	if c.Metrics.Enabled {
		if c.Metrics.Port < 1 || c.Metrics.Port > 65535 {
			return ErrInvalidMetricsPort
		}
	}

	if c.PProf.Enabled {
		if c.PProf.Port < 1 || c.PProf.Port > 65535 {
			return ErrInvalidPProfPort
		}
	}

	if c.HTTP.Address != "" && net.ParseIP(c.HTTP.Address) == nil {
		// If the address is not empty, it must be a valid IP address
		return ErrInvalidHTTPAddress
	}

	if c.Metrics.Enabled && c.Metrics.Address == "" {
		return ErrInvalidMetricsAddress
	}
	if c.PProf.Enabled && c.PProf.Address == "" {
		return ErrInvalidPProfAddress
	}

	for _, proxy := range c.HTTP.TrustedProxies {
		if net.ParseIP(proxy) != nil {
			continue
		}
		if _, _, err := net.ParseCIDR(proxy); err != nil {
			return fmt.Errorf("%w: %q", ErrInvalidTrustedProxies, proxy)
		}
	}

	if c.Storage.Type != StorageTypeMySQL && c.Storage.Type != StorageTypePostgres && c.Storage.Type != StorageTypeSQLite {
		return ErrInvalidStorageType
	}

	if c.Storage.DSN == "" {
		return ErrInvalidStorageDSN
	}

	if c.Auth.JWTSecret == "" {
		return ErrInvalidAuthJWTSecret
	}

	if c.HTTP.URL == "" {
		return ErrInvalidHTTPURL
	}

	return nil
}
