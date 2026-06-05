// Package credstore manages named credential profiles for the groundcover CLI.
//
// A profile bundles the non-secret connection metadata (backend ID, base URL,
// tenant UUID) in a YAML file, while the API key itself lives in the OS keyring.
// This mirrors the design of schpet/linear-cli: the on-disk file never contains
// secrets, only a list of profiles and which one is the default.
package credstore

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

// Profile is the non-secret connection metadata for one named credential.
type Profile struct {
	BackendID  string `yaml:"backend_id,omitempty"`
	BaseURL    string `yaml:"base_url,omitempty"`
	TenantUUID string `yaml:"tenant_uuid,omitempty"`
}

// Store is the on-disk profile metadata. API keys are NOT stored here; they live
// in the OS keyring keyed by profile name.
type Store struct {
	Default  string             `yaml:"default,omitempty"`
	Profiles map[string]Profile `yaml:"profiles"`
}

// Path returns the profiles file location, following the same rules Go's
// os.UserConfigDir uses (XDG_CONFIG_HOME or ~/.config on Unix, %AppData% on
// Windows). The GROUNDCOVER_CONFIG_DIR env var overrides it (used by tests).
func Path() (string, error) {
	if dir := os.Getenv("GROUNDCOVER_CONFIG_DIR"); dir != "" {
		return filepath.Join(dir, "profiles.yaml"), nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("determine config dir: %w", err)
	}
	return filepath.Join(dir, "groundcover", "profiles.yaml"), nil
}

// Load reads the profiles file. A missing file yields an empty store, not an error.
func Load() (*Store, error) {
	path, err := Path()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Store{Profiles: map[string]Profile{}}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read profiles file %s: %w", path, err)
	}
	var s Store
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parse profiles file %s (delete it and re-run `groundcover auth login`): %w", path, err)
	}
	if s.Profiles == nil {
		s.Profiles = map[string]Profile{}
	}
	return &s, nil
}

// Save writes the profiles file, creating parent directories as needed.
func (s *Store) Save() error {
	path, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshal profiles: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write profiles file %s: %w", path, err)
	}
	return nil
}

// Names returns the configured profile names in sorted order.
func (s *Store) Names() []string {
	names := make([]string, 0, len(s.Profiles))
	for name := range s.Profiles {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Add stores the API key in the keyring and records the profile metadata.
// The first profile added automatically becomes the default.
func (s *Store) Add(name string, p Profile, apiKey string) error {
	if name == "" {
		return errors.New("profile name is required")
	}
	if err := active.Set(name, apiKey); err != nil {
		return fmt.Errorf("store API key in keyring for profile %q: %w", name, err)
	}
	if s.Profiles == nil {
		s.Profiles = map[string]Profile{}
	}
	first := len(s.Profiles) == 0
	s.Profiles[name] = p
	if first {
		s.Default = name
	}
	return s.Save()
}

// Remove deletes the profile's keyring entry and metadata. If it was the default,
// the default is reassigned to another profile (or cleared).
func (s *Store) Remove(name string) error {
	if _, ok := s.Profiles[name]; !ok {
		return fmt.Errorf("profile %q not found", name)
	}
	if err := active.Delete(name); err != nil && !errors.Is(err, ErrNotFound) {
		return fmt.Errorf("remove API key from keyring for profile %q: %w", name, err)
	}
	delete(s.Profiles, name)
	if s.Default == name {
		s.Default = ""
		if names := s.Names(); len(names) > 0 {
			s.Default = names[0]
		}
	}
	return s.Save()
}

// SetDefault marks an existing profile as the default.
func (s *Store) SetDefault(name string) error {
	if _, ok := s.Profiles[name]; !ok {
		return fmt.Errorf("profile %q not found", name)
	}
	s.Default = name
	return s.Save()
}

// APIKey returns the keyring-stored API key for a profile.
func (s *Store) APIKey(name string) (string, error) {
	return active.Get(name)
}
