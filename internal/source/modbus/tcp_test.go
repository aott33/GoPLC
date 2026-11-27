package modbus

import (
	"strings"
	"testing"
	"time"

	"github.com/aott33/go-plc/internal/source"
	"github.com/aott33/go-plc/internal/types"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

// configNode parses a YAML mapping and returns its root node.
func configNode(t *testing.T, src string) ast.Node {
	t.Helper()
	f, err := parser.ParseBytes([]byte(src), 0)
	if err != nil {
		t.Fatalf("parse yaml: %v", err)
	}
	return f.Docs[0].Body
}

func TestParseTCPConfigDefaults(t *testing.T) {
	cases := []struct {
		name string
		yaml string // empty means a nil node
	}{
		{"nil node", ""},
		{"null node", "~"},
		{"host only", "host: 10.0.0.5"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var node ast.Node
			if tc.yaml != "" {
				node = configNode(t, tc.yaml)
			}
			sc, err := parseTCPConfig("dev", node)
			if err != nil {
				t.Fatalf("parseTCPConfig: %v", err)
			}
			cfg := sc.(*TCPConfig)

			if cfg.SourceName() != "dev" || cfg.SourceType() != TypeName {
				t.Errorf("name/type = %q/%q", cfg.SourceName(), cfg.SourceType())
			}
			if cfg.Port != DefaultPort {
				t.Errorf("Port = %d, want %d", cfg.Port, DefaultPort)
			}
			if cfg.UnitID != DefaultUnitID {
				t.Errorf("UnitID = %d, want %d", cfg.UnitID, DefaultUnitID)
			}
			if cfg.Timeout.Duration() != DefaultTimeout {
				t.Errorf("Timeout = %v, want %v", cfg.Timeout.Duration(), DefaultTimeout)
			}
			if cfg.PollInterval.Duration() != DefaultPollInterval {
				t.Errorf("PollInterval = %v, want %v", cfg.PollInterval.Duration(), DefaultPollInterval)
			}
			if cfg.RetryInterval.Duration() != DefaultRetryInterval {
				t.Errorf("RetryInterval = %v, want %v", cfg.RetryInterval.Duration(), DefaultRetryInterval)
			}
			if cfg.ByteOrder != DefaultByteOrder || cfg.WordOrder != DefaultWordOrder {
				t.Errorf("orders = %q/%q", cfg.ByteOrder, cfg.WordOrder)
			}
		})
	}
}

func TestParseTCPConfigOverrides(t *testing.T) {
	node := configNode(t, `
host: plc.local
port: 5020
unitId: 0
timeout: 250ms
pollInterval: 100ms
retryInterval: 3s
byteOrder: little-endian
wordOrder: low-word-first
`)
	sc, err := parseTCPConfig("dev", node)
	if err != nil {
		t.Fatalf("parseTCPConfig: %v", err)
	}
	got := *sc.(*TCPConfig)
	want := TCPConfig{
		name:          "dev",
		Host:          "plc.local",
		Port:          5020,
		UnitID:        0, // explicit zero must override the default
		Timeout:       types.Duration(250 * time.Millisecond),
		PollInterval:  types.Duration(100 * time.Millisecond),
		RetryInterval: types.Duration(3 * time.Second),
		ByteOrder:     ByteOrderLittleEndian,
		WordOrder:     WordOrderLowFirst,
	}
	if got != want {
		t.Errorf("got %+v\nwant %+v", got, want)
	}
	if err := got.Validate(); err != nil {
		t.Errorf("Validate: %v", err)
	}
}

func TestParseTCPConfigErrors(t *testing.T) {
	cases := []struct {
		name string
		yaml string
		want string
	}{
		{"unknown field", "host: a\npollIntervall: 1s", "pollIntervall"},
		{"bad duration", "host: a\ntimeout: soon", `[2:10] time: invalid duration "soon"`},
		{"wrong type", "host: a\nport: abc", "port"},
		{"unit id overflow", "host: a\nunitId: 300", "unitId"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseTCPConfig("dev", configNode(t, tc.yaml))
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tc.want) || !strings.Contains(err.Error(), "source: dev") {
				t.Errorf("error %q does not mention %q and the source", err, tc.want)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	valid := func() TCPConfig {
		c := defaultTCPConfig("dev")
		c.Host = "127.0.0.1"
		return c
	}

	cases := []struct {
		name   string
		modify func(*TCPConfig)
		want   string // substring of the error; empty means valid
	}{
		{"valid", func(*TCPConfig) {}, ""},
		{"port 1", func(c *TCPConfig) { c.Port = 1 }, ""},
		{"port 65535", func(c *TCPConfig) { c.Port = 65535 }, ""},
		{"missing host", func(c *TCPConfig) { c.Host = "" }, "Missing required field 'host'"},
		{"port 0", func(c *TCPConfig) { c.Port = 0 }, "'port'"},
		{"port 65536", func(c *TCPConfig) { c.Port = 65536 }, "'port'"},
		{"negative port", func(c *TCPConfig) { c.Port = -1 }, "'port'"},
		{"zero timeout", func(c *TCPConfig) { c.Timeout = 0 }, "'timeout'"},
		{"negative timeout", func(c *TCPConfig) { c.Timeout = types.Duration(-time.Second) }, "'timeout'"},
		{"zero pollInterval", func(c *TCPConfig) { c.PollInterval = 0 }, "'pollInterval'"},
		{"zero retryInterval", func(c *TCPConfig) { c.RetryInterval = 0 }, "'retryInterval'"},
		{"bad byteOrder", func(c *TCPConfig) { c.ByteOrder = "middle-endian" }, "'byteOrder'"},
		{"empty byteOrder", func(c *TCPConfig) { c.ByteOrder = "" }, "'byteOrder'"},
		{"bad wordOrder", func(c *TCPConfig) { c.WordOrder = "ABCD" }, "'wordOrder'"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := valid()
			tc.modify(&cfg)
			err := cfg.Validate()
			if tc.want == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error containing %q", tc.want)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("error %q does not contain %q", err, tc.want)
			}
		})
	}
}

func TestValidateReportsAllErrors(t *testing.T) {
	cfg := TCPConfig{name: "dev"}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error")
	}
	for _, field := range []string{"host", "port", "timeout", "pollInterval", "retryInterval", "byteOrder", "wordOrder"} {
		if !strings.Contains(err.Error(), "'"+field+"'") {
			t.Errorf("error does not mention %q: %v", field, err)
		}
	}
}

func TestRegistered(t *testing.T) {
	sc, err := source.ParseConfig(TypeName, "dev", nil)
	if err != nil {
		t.Fatalf("ParseConfig: %v", err)
	}
	if _, ok := sc.(*TCPConfig); !ok {
		t.Fatalf("ParseConfig returned %T, want *TCPConfig", sc)
	}
}
