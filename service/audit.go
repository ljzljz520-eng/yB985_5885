package service

import (
	"errors"
	"fmt"
	"storeinspection/domain"
)

func (s *Service) RecordEvent(e domain.Event) error {
	if s.DB == nil {
		return errors.New("storage unavailable")
	}
	return s.DB.AppendEvent(e)
}
func (s *Service) AuditTrail(id string) ([]domain.Event, error) { return s.DB.Events(id) }
func (s *Service) Transition(id, to, actor string) error {
	i, e := s.DB.GetInspection(id)
	if e != nil {
		return e
	}
	if !domain.CanTransition(i.Status, to) {
		return fmt.Errorf("cannot transition %s to %s", i.Status, to)
	}
	old := i.Status
	i.Status = to
	if e = s.DB.SaveInspection(i); e != nil {
		return e
	}
	return s.RecordEvent(domain.NewEvent(id+"-"+to, id, "status_changed", actor, old+"->"+to, s.Clock.Now()))
}
func (s *Service) Reopen(id, actor string) error {
	i, e := s.DB.GetInspection(id)
	if e != nil {
		return e
	}
	if i.Status != domain.StatusClosed {
		return errors.New("only closed inspections reopen")
	}
	i.Status = domain.StatusOpen
	i.ClosedAt = i.ClosedAt.AddDate(0, 0, 0)
	if e = s.DB.SaveInspection(i); e != nil {
		return e
	}
	return s.RecordEvent(domain.NewEvent(id+"-reopen", id, "reopened", actor, "", s.Clock.Now()))
}
