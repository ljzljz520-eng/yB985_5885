package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type Snapshot struct {
	Inspection Inspection
	CapturedAt time.Time
	Version    int
}

func EncodeInspection(i Inspection) (string, error) { b, e := json.Marshal(i); return string(b), e }
func DecodeInspection(raw string) (Inspection, error) {
	var i Inspection
	if strings.TrimSpace(raw) == "" {
		return i, errors.New("empty inspection")
	}
	e := json.Unmarshal([]byte(raw), &i)
	return i, e
}
func EncodeFinding(f Finding) (string, error) { b, e := json.Marshal(f); return string(b), e }
func DecodeFinding(raw string) (Finding, error) {
	var f Finding
	e := json.Unmarshal([]byte(raw), &f)
	return f, e
}
func EncodeSnapshot(s Snapshot) (string, error) { b, e := json.Marshal(s); return string(b), e }
func DecodeSnapshot(raw string) (Snapshot, error) {
	var s Snapshot
	e := json.Unmarshal([]byte(raw), &s)
	return s, e
}
func MergeFinding(old, new Finding) Finding {
	out := old
	if new.Title != "" {
		out.Title = new.Title
	}
	if new.Description != "" {
		out.Description = new.Description
	}
	if new.Severity != "" {
		out.Severity = new.Severity
	}
	if new.OwnerID != "" {
		out.OwnerID = new.OwnerID
	}
	if new.Status != "" {
		out.Status = new.Status
	}
	if new.Resolution != "" {
		out.Resolution = new.Resolution
	}
	if !new.DueAt.IsZero() {
		out.DueAt = new.DueAt
	}
	return out
}
func IsStale(i Inspection, now time.Time, limit time.Duration) bool {
	return !i.OpenedAt.IsZero() && now.Sub(i.OpenedAt) > limit && i.Status != StatusClosed
}
