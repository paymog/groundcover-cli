package config

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	sdktransport "github.com/groundcover-com/groundcover-sdk-go/pkg/transport"
	"github.com/paymog/groundcover-cli/internal/credstore"
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
	Profile    string
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

// Resolve fills credentials following this precedence (highest first):
//
//  1. explicit --api-key flag or GROUNDCOVER_API_KEY/GC_API_KEY env var
//  2. --profile flag       -> stored profile (keyring + metadata)
//  3. default profile      -> stored profile (keyring + metadata)
//
// An explicit key combined with --profile is rejected as ambiguous. Explicitly
// set fields (flags/env for backend-id, base-url, tenant-uuid) always win over
// profile-supplied values; the profile only fills gaps. After profile lookup,
// ApplyEnv backfills any remaining unset fields from the environment/defaults.
func (c *Config) Resolve(store *credstore.Store) error {
	envKey := firstEnv("GROUNDCOVER_API_KEY", "GC_API_KEY")
	explicitKey := strings.TrimSpace(c.APIKey) != "" || envKey != ""

	if explicitKey && c.Profile != "" {
		return errors.New("cannot combine --profile with an explicit --api-key/GROUNDCOVER_API_KEY; choose one")
	}

	if !explicitKey {
		name := c.Profile
		if name == "" {
			name = store.Default
		}
		if name != "" {
			p, ok := store.Profiles[name]
			if !ok {
				return fmt.Errorf("profile %q not found; run `groundcover auth login` or `groundcover auth list`", name)
			}
			key, err := store.APIKey(name)
			if err != nil {
				return fmt.Errorf("no keyring entry for profile %q; re-run `groundcover auth login`: %w", name, err)
			}
			c.APIKey = key
			if c.BackendID == "" {
				c.BackendID = p.BackendID
			}
			if c.TenantUUID == "" {
				c.TenantUUID = p.TenantUUID
			}
			if (strings.TrimSpace(c.BaseURL) == "" || c.BaseURL == DefaultBaseURL) && p.BaseURL != "" {
				c.BaseURL = p.BaseURL
			}
		}
	}

	c.ApplyEnv()
	return nil
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
