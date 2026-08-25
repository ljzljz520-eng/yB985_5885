package domain

import (
	"errors"
	"sort"
	"time"
)

type EscalationRule struct {
	Severity   string
	After      time.Duration
	TargetRole string
}
type Escalation struct {
	FindingID, TargetRole, Reason string
	At                            time.Time
}

func DefaultEscalations() []EscalationRule {
	return []EscalationRule{{"high", 24 * time.Hour, "regional-manager"}, {"medium", 72 * time.Hour, "store-manager"}, {"low", 7 * 24 * time.Hour, "store-manager"}}
}
func Escalate(f Finding, now time.Time, rules []EscalationRule) (Escalation, bool) {
	if f.Status == StatusClosed || f.DueAt.IsZero() {
		return Escalation{}, false
	}
	age := now.Sub(f.DueAt)
	for _, r := range rules {
		if r.Severity == f.Severity && age >= r.After {
			return Escalation{FindingID: f.ID, TargetRole: r.TargetRole, Reason: "deadline exceeded", At: now}, true
		}
	}
	return Escalation{}, false
}
func ValidateRules(rules []EscalationRule) error {
	seen := map[string]bool{}
	for _, r := range rules {
		if seen[r.Severity] {
			return errors.New("duplicate severity rule")
		}
		seen[r.Severity] = true
		if SeverityWeight(r.Severity) == 0 || r.After <= 0 {
			return errors.New("invalid escalation rule")
		}
	}
	return nil
}
func SortRules(rules []EscalationRule) {
	sort.Slice(rules, func(i, j int) bool { return rules[i].After < rules[j].After })
}
func DueAt(opened time.Time, severity string) time.Time {
	days := 7
	if severity == "high" {
		days = 1
	} else if severity == "medium" {
		days = 3
	}
	return opened.Add(time.Duration(days) * 24 * time.Hour)
}
