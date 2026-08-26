package domain

import (
	"errors"
	"strings"
	"time"
)

type Store struct {
	ID, Name, Region, Manager string
	Active                    bool
	CreatedAt                 time.Time
}
type Inspection struct {
	ID, StoreID, Inspector, Status, Summary string
	OpenedAt, ClosedAt                      time.Time
	Findings                                []Finding
}
type Finding struct {
	ID, InspectionID, Title, Severity, Status, OwnerID, Description string
	DueAt                                                           time.Time
	Resolution                                                      string
}
type Assignment struct {
	ID, FindingID, OwnerID, Note string
	AssignedAt                   time.Time
	Accepted                     bool
}
type Remediation struct {
	ID, FindingID, Author, Text, Attachment string
	SubmittedAt                             time.Time
	Approved                                bool
}
type User struct {
	ID, Name, Role, Email string
	Enabled               bool
}

const (
	StatusOpen      = "open"
	StatusAssigned  = "assigned"
	StatusSubmitted = "submitted"
	StatusClosed    = "closed"
	StatusRejected  = "rejected"
)

func NormalizeStatus(s string) string { return strings.ToLower(strings.TrimSpace(s)) }
func ValidInspectionStatus(s string) bool {
	switch NormalizeStatus(s) {
	case StatusOpen, StatusAssigned, StatusSubmitted, StatusClosed:
		return true
	}
	return false
}
func ValidateStore(s Store) error {
	if strings.TrimSpace(s.ID) == "" {
		return errors.New("store id required")
	}
	if strings.TrimSpace(s.Name) == "" {
		return errors.New("store name required")
	}
	return nil
}
func ValidateFinding(f Finding) error {
	if f.ID == "" || f.InspectionID == "" {
		return errors.New("finding identifiers required")
	}
	if f.Title == "" {
		return errors.New("finding title required")
	}
	if f.Severity != "low" && f.Severity != "medium" && f.Severity != "high" {
		return errors.New("invalid severity")
	}
	return nil
}
func CanTransition(from, to string) bool {
	from = NormalizeStatus(from)
	to = NormalizeStatus(to)
	if from == StatusOpen && to == StatusAssigned {
		return true
	}
	if from == StatusAssigned && to == StatusSubmitted {
		return true
	}
	if from == StatusSubmitted && (to == StatusClosed || to == StatusRejected) {
		return true
	}
	if from == StatusRejected && to == StatusSubmitted {
		return true
	}
	return false
}
func CloneInspection(in Inspection) Inspection {
	out := in
	out.Findings = append([]Finding(nil), in.Findings...)
	return out
}
