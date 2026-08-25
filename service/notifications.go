package service

import (
	"fmt"
	"storeinspection/domain"
)

type Notification struct{ Recipient, Subject, Body string }
type NotificationSink interface{ Send(Notification) error }

func BuildAssignmentNotice(f domain.Finding, owner string) Notification {
	return Notification{Recipient: owner, Subject: "新的巡店整改任务", Body: fmt.Sprintf("问题 %s：%s，严重级别 %s", f.ID, f.Title, f.Severity)}
}
func BuildSubmissionNotice(f domain.Finding, reviewer string) Notification {
	return Notification{Recipient: reviewer, Subject: "整改说明待审核", Body: fmt.Sprintf("问题 %s 已提交整改说明", f.ID)}
}
func SendNotice(sink NotificationSink, n Notification) error {
	if sink == nil {
		return fmt.Errorf("notification sink unavailable")
	}
	return sink.Send(n)
}
func BuildEscalationNotice(e domain.Escalation) Notification {
	return Notification{Recipient: e.TargetRole, Subject: "巡店问题升级提醒", Body: e.FindingID + ": " + e.Reason}
}
func (s *Service) NotifyAssignment(sink NotificationSink, f domain.Finding, owner string) error {
	return SendNotice(sink, BuildAssignmentNotice(f, owner))
}
