package main

import (
	"errors"
	"flag"
	"io"
	"log/slog"
	"testing"
)

func TestParseFlags(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		want    options
		wantErr bool
	}{
		{"defaults", nil, options{configPath: DefaultConfigPath}, false},
		{"config", []string{"-config", "x.yaml"}, options{configPath: "x.yaml"}, false},
		{"log level", []string{"-log-level", "debug"}, options{configPath: DefaultConfigPath, logLevel: "debug"}, false},
		{"bad log level", []string{"-log-level", "loud"}, options{}, true},
		{"unknown flag", []string{"-nope"}, options{}, true},
		{"positional arg", []string{"extra"}, options{}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseFlags(tc.args, io.Discard)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %+v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestParseFlagsHelp(t *testing.T) {
	_, err := parseFlags([]string{"-h"}, io.Discard)
	if !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("err = %v, want flag.ErrHelp", err)
	}
}

func TestParseLevel(t *testing.T) {
	cases := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
	for name, want := range cases {
		got, err := parseLevel(name)
		if err != nil || got != want {
			t.Errorf("parseLevel(%q) = %v, %v; want %v", name, got, err, want)
		}
	}
	for _, bad := range []string{"", "INFO", "trace"} {
		if _, err := parseLevel(bad); err == nil {
			t.Errorf("parseLevel(%q): expected error", bad)
		}
	}
}
