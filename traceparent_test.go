package utils

import (
	"testing"
)

func TestParseTraceParent(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      *TraceParent
		shouldErr bool
	}{
		{
			name:  "Valid input",
			input: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
			want: &TraceParent{
				Version:    0x00,
				TraceID:    "4bf92f3577b34da6a3ce929d0e0e4736",
				ParentID:   "00f067aa0ba902b7",
				TraceFlags: 0x01,
			},
			shouldErr: false,
		},
		{
			name:      "Invalid format",
			input:     "00-4bf92f3577b34da6a3ce929d0e0e4736-01",
			shouldErr: true,
		},
		{
			name:      "Zero trace-id",
			input:     "00-00000000000000000000000000000000-00f067aa0ba902b7-01",
			shouldErr: true,
		},
		{
			name:      "Zero parent-id",
			input:     "00-4bf92f3577b34da6a3ce929d0e0e4736-0000000000000000-01",
			shouldErr: true,
		},
		{
			name:      "Invalid hex in version",
			input:     "zz-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
			shouldErr: true,
		},
		{
			name:      "Invalid field length",
			input:     "00-1234-00f067aa0ba902b7-01",
			shouldErr: true,
		},
		{
			name:      "Invalid flags",
			input:     "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-gh",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTraceParent(tt.input)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("expected error, got nil for input: %q", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v for input: %q", err, tt.input)
				}
				if got.Version != tt.want.Version ||
					got.TraceID != tt.want.TraceID ||
					got.ParentID != tt.want.ParentID ||
					got.TraceFlags != tt.want.TraceFlags {
					t.Errorf("got %+v, want %+v", got, tt.want)
				}
			}
		})
	}
}
