package source

import (
	"slices"
	"strings"
	"testing"

	"github.com/goccy/go-yaml/ast"
)

type stubConfig struct{ name string }

func (s *stubConfig) Validate() error    { return nil }
func (s *stubConfig) SourceType() string { return "stub" }
func (s *stubConfig) SourceName() string { return s.name }

func stubFactory(name string, _ ast.Node) (SourceConfig, error) {
	return &stubConfig{name: name}, nil
}

// withRegistry swaps in an empty registry for the duration of a test.
func withRegistry(t *testing.T) {
	t.Helper()
	mu.Lock()
	saved := registry
	registry = make(map[string]SourceFactory)
	mu.Unlock()
	t.Cleanup(func() {
		mu.Lock()
		registry = saved
		mu.Unlock()
	})
}

func TestRegisteredTypesSorted(t *testing.T) {
	withRegistry(t)
	Register("zeta", stubFactory)
	Register("alpha", stubFactory)
	Register("mid", stubFactory)

	got := RegisteredTypes()
	want := []string{"alpha", "mid", "zeta"}
	if !slices.Equal(got, want) {
		t.Fatalf("RegisteredTypes() = %v, want %v", got, want)
	}
}

func TestParseConfig(t *testing.T) {
	withRegistry(t)
	Register("zeta", stubFactory)
	Register("alpha", stubFactory)

	cfg, err := ParseConfig("alpha", "dev1", nil)
	if err != nil {
		t.Fatalf("ParseConfig: unexpected error: %v", err)
	}
	if cfg.SourceName() != "dev1" {
		t.Errorf("SourceName() = %q, want %q", cfg.SourceName(), "dev1")
	}

	_, err = ParseConfig("nope", "dev2", nil)
	if err == nil {
		t.Fatal("ParseConfig with unknown type: expected error")
	}
	msg := err.Error()
	for _, want := range []string{"'nope'", "source: dev2", "registered: alpha, zeta"} {
		if !strings.Contains(msg, want) {
			t.Errorf("error %q does not contain %q", msg, want)
		}
	}
}

func TestRegisterDuplicatePanics(t *testing.T) {
	withRegistry(t)
	Register("alpha", stubFactory)

	defer func() {
		if recover() == nil {
			t.Fatal("Register with duplicate type: expected panic")
		}
	}()
	Register("alpha", stubFactory)
}
