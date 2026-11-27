package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aott33/go-plc/internal/source/modbus"
)

func TestLoadExample(t *testing.T) {
	cfg, err := Load(filepath.Join("..", "..", "configs", "example.yaml"))
	if err != nil {
		t.Fatalf("Load example: %v", err)
	}
	if len(cfg.Sources) != 1 {
		t.Fatalf("got %d sources, want 1", len(cfg.Sources))
	}
	src := cfg.Sources[0]
	if src.Name != "simulator" || src.Type != modbus.TypeName {
		t.Errorf("source = %q/%q", src.Name, src.Type)
	}
	tcp, ok := src.Config.(*modbus.TCPConfig)
	if !ok {
		t.Fatalf("Config is %T, want *modbus.TCPConfig", src.Config)
	}
	if tcp.Host != "127.0.0.1" || tcp.Port != 5020 {
		t.Errorf("host:port = %s:%d", tcp.Host, tcp.Port)
	}
}

func TestLoadDefaults(t *testing.T) {
	cfg, err := parse("test.yaml", []byte(`
sources:
  - name: plc1
    type: modbus-tcp
    config:
      host: 10.0.0.1
`))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if cfg.LogLevel != DefaultLogLevel {
		t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, DefaultLogLevel)
	}
	tcp := cfg.Sources[0].Config.(*modbus.TCPConfig)
	if tcp.Port != modbus.DefaultPort ||
		tcp.UnitID != modbus.DefaultUnitID ||
		tcp.Timeout.Duration() != modbus.DefaultTimeout ||
		tcp.PollInterval.Duration() != modbus.DefaultPollInterval ||
		tcp.RetryInterval.Duration() != modbus.DefaultRetryInterval ||
		tcp.ByteOrder != modbus.DefaultByteOrder ||
		tcp.WordOrder != modbus.DefaultWordOrder {
		t.Errorf("defaults not applied: %+v", *tcp)
	}
}

func TestLoadValues(t *testing.T) {
	cfg, err := parse("test.yaml", []byte(`
logLevel: debug
sources:
  - name: a
    type: modbus-tcp
    config: {host: h1, pollInterval: 250ms}
  - name: b
    type: modbus-tcp
    config: {host: h2, port: 1502}
variables:
  - name: tankLevel
    source: a
    dataType: uint16
`))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("LogLevel = %q", cfg.LogLevel)
	}
	if len(cfg.Sources) != 2 || cfg.Sources[0].Name != "a" || cfg.Sources[1].Name != "b" {
		t.Fatalf("sources = %+v", cfg.Sources)
	}
	if got := cfg.Sources[0].Config.(*modbus.TCPConfig).PollInterval.Duration(); got != 250*time.Millisecond {
		t.Errorf("a.pollInterval = %v", got)
	}
	if got := cfg.Sources[1].Config.(*modbus.TCPConfig).Port; got != 1502 {
		t.Errorf("b.port = %d", got)
	}
	if len(cfg.Variables) != 1 || cfg.Variables[0].Name != "tankLevel" {
		t.Errorf("variables = %+v", cfg.Variables)
	}
}

func TestLoadErrors(t *testing.T) {
	cases := []struct {
		name string
		yaml string
		want []string
	}{
		{
			name: "unknown type lists registered types",
			yaml: "sources:\n  - {name: x, type: opcua-client}\n",
			want: []string{"Unknown source type 'opcua-client'", "source: x", "registered: modbus-tcp"},
		},
		{
			name: "missing host",
			yaml: "sources:\n  - {name: x, type: modbus-tcp, config: {port: 502}}\n",
			want: []string{"Missing required field 'host'", "source: x"},
		},
		{
			name: "missing config block",
			yaml: "sources:\n  - {name: x, type: modbus-tcp}\n",
			want: []string{"'host'"},
		},
		{
			name: "port out of range",
			yaml: "sources:\n  - {name: x, type: modbus-tcp, config: {host: h, port: 70000}}\n",
			want: []string{"'port'", "value: 70000"},
		},
		{
			name: "zero timeout",
			yaml: "sources:\n  - {name: x, type: modbus-tcp, config: {host: h, timeout: 0s}}\n",
			want: []string{"'timeout'"},
		},
		{
			name: "zero pollInterval",
			yaml: "sources:\n  - {name: x, type: modbus-tcp, config: {host: h, pollInterval: 0s}}\n",
			want: []string{"'pollInterval'"},
		},
		{
			name: "zero retryInterval",
			yaml: "sources:\n  - {name: x, type: modbus-tcp, config: {host: h, retryInterval: 0s}}\n",
			want: []string{"'retryInterval'"},
		},
		{
			name: "bad byteOrder",
			yaml: "sources:\n  - {name: x, type: modbus-tcp, config: {host: h, byteOrder: BE}}\n",
			want: []string{"'byteOrder'"},
		},
		{
			name: "bad wordOrder",
			yaml: "sources:\n  - {name: x, type: modbus-tcp, config: {host: h, wordOrder: swapped}}\n",
			want: []string{"'wordOrder'"},
		},
		{
			name: "unknown field in source config",
			yaml: "sources:\n  - {name: x, type: modbus-tcp, config: {host: h, hots: h}}\n",
			want: []string{"hots", "source: x"},
		},
		{
			name: "missing name",
			yaml: "sources:\n  - {type: modbus-tcp}\n",
			want: []string{"Missing required field 'name'", "source index: 0"},
		},
		{
			name: "missing type",
			yaml: "sources:\n  - {name: x}\n",
			want: []string{"Missing required field 'type'"},
		},
		{
			name: "duplicate name",
			yaml: "sources:\n  - {name: x, type: modbus-tcp, config: {host: h}}\n  - {name: x, type: modbus-tcp, config: {host: h}}\n",
			want: []string{"Duplicate source name", "source: x"},
		},
		{
			name: "bad log level",
			yaml: "logLevel: loud\n",
			want: []string{"'logLevel'", "loud"},
		},
		{
			name: "invalid yaml",
			yaml: "sources: [\n",
			want: []string{"Invalid YAML", "path: test.yaml"},
		},
		{
			name: "errors from several sources are all reported",
			yaml: "sources:\n  - {name: a, type: modbus-tcp}\n  - {name: b, type: nope}\n",
			want: []string{"source: a", "source: b"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parse("test.yaml", []byte(tc.yaml))
			if err == nil {
				t.Fatal("expected error")
			}
			for _, w := range tc.want {
				if !strings.Contains(err.Error(), w) {
					t.Errorf("error %q does not contain %q", err, w)
				}
			}
		})
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "missing.yaml"))
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Load missing file: err = %v, want os.ErrNotExist", err)
	}
}

func TestLoadFromFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "c.yaml")
	data := "sources:\n  - {name: x, type: modbus-tcp, config: {host: h}}\n"
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Sources) != 1 {
		t.Fatalf("got %d sources", len(cfg.Sources))
	}
}
