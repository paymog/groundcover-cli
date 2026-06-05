package config_test

import (
	"strings"
	"testing"

	"github.com/paymog/groundcover-cli/internal/config"
	"github.com/paymog/groundcover-cli/internal/credstore"
)

type fakeKeyring struct{ m map[string]string }

func (f fakeKeyring) Get(a string) (string, error) {
	v, ok := f.m[a]
	if !ok {
		return "", credstore.ErrNotFound
	}
	return v, nil
}
func (f fakeKeyring) Set(a, s string) error { f.m[a] = s; return nil }
func (f fakeKeyring) Delete(a string) error { delete(f.m, a); return nil }

func setup(t *testing.T, keys map[string]string) {
	t.Helper()
	prev := credstore.SetBackend(fakeKeyring{m: keys})
	t.Cleanup(func() { credstore.SetBackend(prev) })
	// Isolate env that the resolver reads.
	for _, k := range []string{"GROUNDCOVER_API_KEY", "GC_API_KEY", "GROUNDCOVER_BACKEND_ID", "GC_BACKEND_ID"} {
		t.Setenv(k, "")
	}
}

func TestResolveDefaultProfile(t *testing.T) {
	setup(t, map[string]string{"acme": "key-acme"})
	store := &credstore.Store{
		Default:  "acme",
		Profiles: map[string]credstore.Profile{"acme": {BackendID: "acme-be"}},
	}
	c := config.FromEnv()
	if err := c.Resolve(store); err != nil {
		t.Fatal(err)
	}
	if c.APIKey != "key-acme" || c.BackendID != "acme-be" {
		t.Fatalf("got key=%q backend=%q", c.APIKey, c.BackendID)
	}
}

func TestResolveNamedProfileOverridesDefault(t *testing.T) {
	setup(t, map[string]string{"acme": "k1", "side": "k2"})
	store := &credstore.Store{
		Default: "acme",
		Profiles: map[string]credstore.Profile{
			"acme": {BackendID: "acme-be"},
			"side": {BackendID: "side-be"},
		},
	}
	c := config.FromEnv()
	c.Profile = "side"
	if err := c.Resolve(store); err != nil {
		t.Fatal(err)
	}
	if c.APIKey != "k2" || c.BackendID != "side-be" {
		t.Fatalf("got key=%q backend=%q", c.APIKey, c.BackendID)
	}
}

func TestResolveExplicitKeyBeatsProfile(t *testing.T) {
	setup(t, map[string]string{"acme": "k1"})
	store := &credstore.Store{Default: "acme", Profiles: map[string]credstore.Profile{"acme": {BackendID: "acme-be"}}}
	c := config.FromEnv()
	c.APIKey = "flag-key"
	if err := c.Resolve(store); err != nil {
		t.Fatal(err)
	}
	if c.APIKey != "flag-key" {
		t.Fatalf("explicit key should win, got %q", c.APIKey)
	}
	// No profile selected, so backend falls back to the default backend ID.
	if c.BackendID != config.DefaultBackendID {
		t.Fatalf("backend = %q", c.BackendID)
	}
}

func TestResolveAmbiguousKeyAndProfile(t *testing.T) {
	setup(t, nil)
	c := config.FromEnv()
	c.APIKey = "flag-key"
	c.Profile = "acme"
	err := c.Resolve(&credstore.Store{})
	if err == nil || !strings.Contains(err.Error(), "cannot combine") {
		t.Fatalf("expected ambiguity error, got %v", err)
	}
}

func TestResolveMissingNamedProfile(t *testing.T) {
	setup(t, nil)
	c := config.FromEnv()
	c.Profile = "ghost"
	err := c.Resolve(&credstore.Store{Profiles: map[string]credstore.Profile{}})
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected not-found error, got %v", err)
	}
}

func TestResolveNoProfileNoEnvIsEmpty(t *testing.T) {
	setup(t, nil)
	c := config.FromEnv()
	if err := c.Resolve(&credstore.Store{Profiles: map[string]credstore.Profile{}}); err != nil {
		t.Fatal(err)
	}
	if c.RequireAPIKey() == nil {
		t.Fatal("expected RequireAPIKey to fail when nothing is configured")
	}
}
