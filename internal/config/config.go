package config

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	sdktransport "github.com/groundcover-com/groundcover-sdk-go/pkg/transport"
)

const (
	DefaultBaseURL = "https://api.groundcover.com"
	// DefaultBackendID is groundcover's standard backend ID. It is defaulted so the
	// common case spares every command site from re-passing --backend-id; override
	// with GROUNDCOVER_BACKEND_ID or --backend-id for other backends.
	DefaultBackendID = "groundcover"
)

type Config struct {
	APIKey     string
	BackendID  string
	BaseURL    string
	TenantUUID string
	Raw        bool
	Timeout    time.Duration
}

func FromEnv() Config {
	return Config{
		BaseURL: defaultString(firstEnv("GROUNDCOVER_BASE_URL", "GC_BASE_URL"), DefaultBaseURL),
		Timeout: 30 * time.Second,
	}
}

func (c *Config) ApplyEnv() {
	if c.APIKey == "" {
		c.APIKey = firstEnv("GROUNDCOVER_API_KEY", "GC_API_KEY")
	}
	if c.BackendID == "" {
		c.BackendID = defaultString(firstEnv("GROUNDCOVER_BACKEND_ID", "GC_BACKEND_ID"), DefaultBackendID)
	}
	if c.TenantUUID == "" {
		c.TenantUUID = firstEnv("GROUNDCOVER_TENANT_UUID", "GC_TENANT_UUID")
	}
	if strings.TrimSpace(c.BaseURL) == "" {
		c.BaseURL = defaultString(firstEnv("GROUNDCOVER_BASE_URL", "GC_BASE_URL"), DefaultBaseURL)
	}
}

func (c Config) NormalizedBaseURL() string {
	if strings.TrimSpace(c.BaseURL) == "" {
		return DefaultBaseURL
	}
	return strings.TrimRight(c.BaseURL, "/")
}

func (c Config) RequireAPIKey() error {
	if strings.TrimSpace(c.APIKey) == "" {
		return errors.New("missing API key: set GROUNDCOVER_API_KEY or pass --api-key")
	}
	return nil
}

func (c Config) RequireSDKAuth() error {
	if err := c.RequireAPIKey(); err != nil {
		return err
	}
	if strings.TrimSpace(c.BackendID) == "" {
		return errors.New("missing backend ID: set GROUNDCOVER_BACKEND_ID or pass --backend-id")
	}
	return nil
}

func (c Config) HTTPClient() *http.Client {
	return &http.Client{
		Timeout: c.Timeout,
		Transport: sdktransport.NewTransport(
			c.APIKey,
			c.BackendID,
			http.DefaultTransport,
			0,
			0,
			0,
			nil,
		),
	}
}

func firstEnv(names ...string) string {
	for _, name := range names {
		if value := os.Getenv(name); value != "" {
			return value
		}
	}
	return ""
}

func defaultString(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
