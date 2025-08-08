package utils

import (
	"encoding/hex"
	"strings"
)

type TraceParent struct {
	Version    byte   // 1 byte, as int
	TraceID    string // 32 hex chars, 16 bytes, lower-case
	ParentID   string // 16 hex chars, 8 bytes, lower-case
	TraceFlags byte   // 1 byte, as int
}

// IsTraced returns true if the trace is sampled (traced).
// According to W3C trace context spec, the least significant bit indicates sampling.
func (tp *TraceParent) IsTraced() bool {
	return tp.TraceFlags&0x01 != 0
}

// ParseTraceParent parses a W3C traceparent header string.
// Returns TraceParent or error on invalid format.
// https://www.w3.org/TR/trace-context/#version-format
func ParseTraceParent(header string) (*TraceParent, error) {
	parts := strings.Split(header, "-")
	if len(parts) != 4 {
		return nil, Errorfc("invalid traceparent format")
	}

	verStr, traceID, parentID, flagsStr := parts[0], parts[1], parts[2], parts[3]

	if len(verStr) != 2 || len(traceID) != 32 || len(parentID) != 16 || len(flagsStr) != 2 {
		return nil, Errorfc("invalid traceparent field lengths")
	}

	ver, err := hex.DecodeString(verStr)
	if err != nil {
		return nil, Errorfc("invalid version: %w", err)
	}
	if string(traceID) == "00000000000000000000000000000000" {
		return nil, Errorfc("trace-id all zero")
	}
	if string(parentID) == "0000000000000000" {
		return nil, Errorfc("parent-id all zero")
	}
	flags, err := hex.DecodeString(flagsStr)
	if err != nil {
		return nil, Errorfc("invalid trace-flags: %w", err)
	}
	return &TraceParent{
		Version:    ver[0],
		TraceID:    strings.ToLower(traceID),
		ParentID:   strings.ToLower(parentID),
		TraceFlags: flags[0],
	}, nil
}
