package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func ValidateInspection(i Inspection) error {
	if i.ID == "" {
		return errors.New("inspection id required")
	}
	if i.StoreID == "" {
		return errors.New("store required")
	}
	if i.Inspector == "" {
		return errors.New("inspector required")
	}
	if !ValidInspectionStatus(i.Status) {
		return errors.New("unsupported inspection status")
	}
	return nil
}
func ValidateAssignment(a Assignment) error {
	if a.ID == "" || a.FindingID == "" {
		return errors.New("assignment identifiers required")
	}
	if a.OwnerID == "" {
		return errors.New("assignment owner required")
	}
	return nil
}
func ValidateRemediation(r Remediation) error {
	if r.ID == "" || r.FindingID == "" {
		return errors.New("remediation identifiers required")
	}
	if strings.TrimSpace(r.Text) == "" {
		return errors.New("remediation text required")
	}
	if r.SubmittedAt.IsZero() {
		return errors.New("submission time required")
	}
	return nil
}
func Age(due, now time.Time) string {
	if due.IsZero() {
		return "unscheduled"
	}
	if due.Before(now) {
		return "overdue"
	}
	if due.Sub(now) < 48*time.Hour {
		return "due-soon"
	}
	return "planned"
}
func StatusMessage(status string) string {
	switch NormalizeStatus(status) {
	case StatusOpen:
		return "待分配负责人"
	case StatusAssigned:
		return "等待整改说明"
	case StatusSubmitted:
		return "等待审核关闭"
	case StatusClosed:
		return "已关闭"
	default:
		return fmt.Sprintf("未知状态:%s", status)
	}
}
func NormalizeRole(r string) string { return strings.ToLower(strings.TrimSpace(r)) }
