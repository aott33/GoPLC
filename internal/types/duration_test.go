package types

import (
	"testing"
	"time"

	"github.com/goccy/go-yaml"
)

func TestDuration_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		want    time.Duration
		wantErr bool
	}{
		{
			name: "milliseconds",
			yaml: `timeout: "100ms"`,
			want: 100 * time.Millisecond,
		},
		{
			name: "seconds",
			yaml: `timeout: "5s"`,
			want: 5 * time.Second,
		},
		{
			name: "minutes",
			yaml: `timeout: "2m"`,
			want: 2 * time.Minute,
		},
		{
			name: "hours",
			yaml: `timeout: "1h"`,
			want: 1 * time.Hour,
		},
		{
			name: "composite duration",
			yaml: `timeout: "1h30m"`,
			want: 1*time.Hour + 30*time.Minute,
		},
		{
			name: "microseconds",
			yaml: `timeout: "500us"`,
			want: 500 * time.Microsecond,
		},
		{
			name: "nanoseconds",
			yaml: `timeout: "1000ns"`,
			want: 1000 * time.Nanosecond,
		},
		{
			name:    "invalid format - missing unit",
			yaml:    `timeout: "100"`,
			wantErr: true,
		},
		{
			name:    "invalid format - bad unit",
			yaml:    `timeout: "100days"`,
			wantErr: true,
		},
		{
			name:    "invalid format - empty string",
			yaml:    `timeout: ""`,
			wantErr: true,
		},
		{
			name:    "invalid format - not a string",
			yaml:    `timeout: 100`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config struct {
				Timeout Duration `yaml:"timeout"`
			}

			err := yaml.Unmarshal([]byte(tt.yaml), &config)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalYAML() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("UnmarshalYAML() unexpected error: %v", err)
				return
			}

			got := config.Timeout.Duration()
			if got != tt.want {
				t.Errorf("UnmarshalYAML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDuration_Duration(t *testing.T) {
	tests := []struct {
		name     string
		duration Duration
		want     time.Duration
	}{
		{
			name:     "zero duration",
			duration: Duration(0),
			want:     0,
		},
		{
			name:     "positive duration",
			duration: Duration(5 * time.Second),
			want:     5 * time.Second,
		},
		{
			name:     "large duration",
			duration: Duration(24 * time.Hour),
			want:     24 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.duration.Duration()
			if got != tt.want {
				t.Errorf("Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestDuration_RoundTrip tests that we can marshal and unmarshal correctly
func TestDuration_RoundTrip(t *testing.T) {
	type Config struct {
		Timeout       Duration `yaml:"timeout"`
		PollInterval  Duration `yaml:"pollInterval"`
		RetryInterval Duration `yaml:"retryInterval"`
	}

	original := Config{
		Timeout:       Duration(5 * time.Second),
		PollInterval:  Duration(100 * time.Millisecond),
		RetryInterval: Duration(30 * time.Second),
	}

	// Marshal to YAML
	data, err := yaml.Marshal(&original)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}

	// Unmarshal back
	var decoded Config
	if err := yaml.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal() error: %v", err)
	}

	// Compare
	if decoded.Timeout.Duration() != original.Timeout.Duration() {
		t.Errorf("Timeout: got %v, want %v", decoded.Timeout.Duration(), original.Timeout.Duration())
	}
	if decoded.PollInterval.Duration() != original.PollInterval.Duration() {
		t.Errorf("PollInterval: got %v, want %v", decoded.PollInterval.Duration(), original.PollInterval.Duration())
	}
	if decoded.RetryInterval.Duration() != original.RetryInterval.Duration() {
		t.Errorf("RetryInterval: got %v, want %v", decoded.RetryInterval.Duration(), original.RetryInterval.Duration())
	}
}
