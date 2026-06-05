package raw

import "testing"

func TestSetDeep(t *testing.T) {
	target := map[string]any{}
	setDeep(target, "a.b.c", 1)

	a, ok := target["a"].(map[string]any)
	if !ok {
		t.Fatalf("expected nested map at a")
	}
	b, ok := a["b"].(map[string]any)
	if !ok {
		t.Fatalf("expected nested map at a.b")
	}
	if b["c"] != 1 {
		t.Fatalf("expected c=1, got %#v", b["c"])
	}
}

func TestFind(t *testing.T) {
	command, ok := Find([]string{"backend", "settings"})
	if !ok {
		t.Fatalf("expected backend settings command")
	}
	if command.Path != "/api/backend/settings" {
		t.Fatalf("unexpected path %s", command.Path)
	}
}
