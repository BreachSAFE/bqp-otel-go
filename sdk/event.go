// SPDX-FileCopyrightText: 2026 BreachSAFE <https://www.breachsafe.io>
// SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0

// Package sdk provides a bounded, fail-open producer boundary for BQP lifecycle events.
package sdk

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

const (
	Schema             = "bqp.run.v1"
	MaxEventLength     = 128
	MaxComponentLength = 128
	MaxIDLength        = 128
	MaxAttributes      = 32
	MaxValueLength     = 512
)

var (
	ErrInvalidEvent = errors.New("invalid BQP event")
	ErrSensitiveKey = errors.New("sensitive attribute key")
)

// Outcome is the terminal state of a lifecycle event.
type Outcome string

const (
	OutcomeOK      Outcome = "ok"
	OutcomeError   Outcome = "error"
	OutcomeTimeout Outcome = "timeout"
)

// Event is the language-neutral BQP envelope carried as OTel span attributes/events.
// Attributes are copied by NewEvent and must contain only bounded, non-sensitive metadata.
type Event struct {
	Event      string            `json:"event"`
	Component  string            `json:"component"`
	Outcome    Outcome           `json:"outcome"`
	RunID      string            `json:"run_id,omitempty"`
	SessionID  string            `json:"session_id,omitempty"`
	DurationMS int64             `json:"duration_ms,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// NewEvent validates and copies an event before it crosses the telemetry boundary.
func NewEvent(event, component string, outcome Outcome, runID, sessionID string, durationMS int64, attributes map[string]string) (Event, error) {
	result := Event{Event: event, Component: component, Outcome: outcome, RunID: runID, SessionID: sessionID, DurationMS: durationMS, Attributes: attributes}
	if err := result.validate(); err != nil {
		return Event{}, err
	}
	copyAttrs := make(map[string]string, len(attributes))
	for key, value := range attributes {
		copyAttrs[key] = value
	}
	result.Attributes = copyAttrs
	return result, nil
}

func (e Event) validate() error {
	event, component, outcome, runID, sessionID, durationMS, attributes := e.Event, e.Component, e.Outcome, e.RunID, e.SessionID, e.DurationMS, e.Attributes
	if !bounded(event, MaxEventLength) || !bounded(component, MaxComponentLength) || !bounded(runID, MaxIDLength) || !bounded(sessionID, MaxIDLength) {
		return fmt.Errorf("%w: identity exceeds limit", ErrInvalidEvent)
	}
	if event == "" || component == "" {
		return fmt.Errorf("%w: event and component are required", ErrInvalidEvent)
	}
	if outcome != OutcomeOK && outcome != OutcomeError && outcome != OutcomeTimeout {
		return fmt.Errorf("%w: unsupported outcome %q", ErrInvalidEvent, outcome)
	}
	if durationMS < 0 {
		return fmt.Errorf("%w: negative duration", ErrInvalidEvent)
	}
	if len(attributes) > MaxAttributes {
		return fmt.Errorf("%w: too many attributes", ErrInvalidEvent)
	}
	for key, value := range attributes {
		if !bounded(key, MaxEventLength) || key == "" || sensitiveKey(key) {
			if sensitiveKey(key) {
				return fmt.Errorf("%w: %q", ErrSensitiveKey, key)
			}
			return fmt.Errorf("%w: invalid attribute key", ErrInvalidEvent)
		}
		if !bounded(value, MaxValueLength) {
			return fmt.Errorf("%w: attribute %q exceeds limit", ErrInvalidEvent, key)
		}
	}
	return nil
}

// AttributeKeys returns deterministic keys for diagnostics and tests.
func (e Event) AttributeKeys() []string {
	keys := make([]string, 0, len(e.Attributes))
	for key := range e.Attributes {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func bounded(value string, limit int) bool {
	return len(value) <= limit && !strings.ContainsAny(value, "\r\n")
}

func sensitiveKey(key string) bool {
	normalized := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(key, "-", "_"), ".", "_"))
	for _, token := range []string{"token", "secret", "password", "authorization", "oauth", "private_key", "terminal", "stdout", "stderr"} {
		if strings.Contains(normalized, token) {
			return true
		}
	}
	return false
}
