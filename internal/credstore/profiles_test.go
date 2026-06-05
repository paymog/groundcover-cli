package credstore

import (
	"errors"
	"testing"
)

// fakeKeyring is an in-memory Backend for tests.
type fakeKeyring struct{ m map[string]string }

func newFakeKeyring() *fakeKeyring { return &fakeKeyring{m: map[string]string{}} }

func (f *fakeKeyring) Get(account string) (string, error) {
	v, ok := f.m[account]
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}
func (f *fakeKeyring) Set(account, secret string) error { f.m[account] = secret; return nil }
func (f *fakeKeyring) Delete(account string) error {
	if _, ok := f.m[account]; !ok {
		return ErrNotFound
	}
	delete(f.m, account)
	return nil
}

// withFakes installs a fake keyring and a temp config dir for the duration of t.
func withFakes(t *testing.T) *fakeKeyring {
	t.Helper()
	fk := newFakeKeyring()
	prev := SetBackend(fk)
	t.Cleanup(func() { SetBackend(prev) })
	t.Setenv("GROUNDCOVER_CONFIG_DIR", t.TempDir())
	return fk
}

func TestAddFirstProfileBecomesDefault(t *testing.T) {
	fk := withFakes(t)
	s, _ := Load()
	if err := s.Add("acme", Profile{BackendID: "acme-be"}, "key-1"); err != nil {
		t.Fatal(err)
	}
	if s.Default != "acme" {
		t.Fatalf("first profile should be default, got %q", s.Default)
	}
	if fk.m["acme"] != "key-1" {
		t.Fatalf("key not stored in keyring: %q", fk.m["acme"])
	}

	// Second profile must not steal the default.
	if err := s.Add("side", Profile{BackendID: "side-be"}, "key-2"); err != nil {
		t.Fatal(err)
	}
	if s.Default != "acme" {
		t.Fatalf("default changed unexpectedly to %q", s.Default)
	}
}

func TestLoadPersistsMetadataButNotKeys(t *testing.T) {
	withFakes(t)
	s, _ := Load()
	if err := s.Add("acme", Profile{BackendID: "acme-be", BaseURL: "https://x"}, "secret"); err != nil {
		t.Fatal(err)
	}

	// Reload from disk; metadata persists, key comes from keyring.
	reloaded, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.Default != "acme" || reloaded.Profiles["acme"].BackendID != "acme-be" {
		t.Fatalf("metadata not persisted: %+v", reloaded)
	}
	key, err := reloaded.APIKey("acme")
	if err != nil || key != "secret" {
		t.Fatalf("APIKey = %q, %v", key, err)
	}
}

func TestRemoveReassignsDefault(t *testing.T) {
	fk := withFakes(t)
	s, _ := Load()
	_ = s.Add("acme", Profile{BackendID: "a"}, "k1")
	_ = s.Add("side", Profile{BackendID: "b"}, "k2")

	if err := s.Remove("acme"); err != nil {
		t.Fatal(err)
	}
	if s.Default != "side" {
		t.Fatalf("default not reassigned, got %q", s.Default)
	}
	if _, ok := fk.m["acme"]; ok {
		t.Fatal("keyring entry not removed")
	}
	if err := s.Remove("missing"); err == nil {
		t.Fatal("removing missing profile should error")
	}
}

func TestSetDefaultUnknownErrors(t *testing.T) {
	withFakes(t)
	s, _ := Load()
	if err := s.SetDefault("nope"); err == nil {
		t.Fatal("expected error for unknown profile")
	}
}

func TestAPIKeyMissing(t *testing.T) {
	withFakes(t)
	s, _ := Load()
	_, err := s.APIKey("ghost")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected not-found, got %v", err)
	}
}
