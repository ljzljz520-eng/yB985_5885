package service

import (
	"errors"
	"storeinspection/domain"
	"storeinspection/store"
	"time"
)

type Clock interface{ Now() time.Time }
type FixedClock struct{ Value time.Time }

func (c FixedClock) Now() time.Time { return c.Value }

type Service struct {
	DB       *store.DB
	Clock    Clock
	Notifier func(string) error
}

func New(db *store.DB, c Clock) *Service { return &Service{DB: db, Clock: c} }
func (s *Service) OpenInspection(i domain.Inspection) error {
	if s.DB == nil {
		return errors.New("storage unavailable")
	}
	if i.Status == "" {
		i.Status = domain.StatusOpen
	}
	if !domain.ValidInspectionStatus(i.Status) {
		return errors.New("invalid status")
	}
	if i.OpenedAt.IsZero() {
		i.OpenedAt = s.Clock.Now()
	}
	return s.DB.SaveInspection(i)
}
func (s *Service) AddFinding(f domain.Finding) error {
	if e := domain.ValidateFinding(f); e != nil {
		return e
	}
	f.Status = domain.StatusOpen
	return s.DB.SaveFinding(f)
}
func (s *Service) Assign(findingID, owner, note string) error {
	if owner == "" {
		return errors.New("owner required")
	}
	a := domain.Assignment{ID: findingID + "-assignment", FindingID: findingID, OwnerID: owner, Note: note, AssignedAt: s.Clock.Now()}
	if e := s.DB.SaveAssignment(a); e != nil {
		return e
	}
	f, e := s.finding(findingID)
	if e != nil {
		return e
	}
	f.OwnerID = owner
	f.Status = domain.StatusAssigned
	return s.DB.SaveFinding(f)
}
func (s *Service) Submit(findingID, author, text, attachment string) error {
	if text == "" {
		return errors.New("remediation text required")
	}
	f, e := s.finding(findingID)
	if e != nil {
		return e
	}
	if f.Status != domain.StatusAssigned && f.Status != domain.StatusRejected {
		return errors.New("finding is not ready for submission")
	}
	if e = s.DB.SaveRemediation(domain.Remediation{ID: findingID + "-remediation", FindingID: findingID, Author: author, Text: text, Attachment: attachment, SubmittedAt: s.Clock.Now()}); e != nil {
		return e
	}
	f.Status = domain.StatusSubmitted
	if e = s.DB.SaveFinding(f); e != nil {
		return e
	}
	i, e := s.DB.GetInspection(f.InspectionID)
	if e == nil && i.Status != domain.StatusSubmitted {
		i.Status = domain.StatusSubmitted
		e = s.DB.SaveInspection(i)
	}
	return e
}
func (s *Service) Close(inspectionID string) error {
	i, e := s.DB.GetInspection(inspectionID)
	if e != nil {
		return e
	}
	fs, e := s.DB.ListFindings(inspectionID)
	if e != nil {
		return e
	}
	for _, f := range fs {
		if f.Status != domain.StatusSubmitted && f.Status != domain.StatusClosed {
			return errors.New("all findings require remediation")
		}
	}
	if !domain.CanTransition(i.Status, domain.StatusClosed) && i.Status != domain.StatusClosed {
		return errors.New("invalid inspection transition")
	}
	i.Status = domain.StatusClosed
	i.ClosedAt = s.Clock.Now()
	return s.DB.SaveInspection(i)
}
func (s *Service) finding(id string) (domain.Finding, error) {
	i, e := s.DB.ListInspections()
	_ = i
	if e != nil {
		return domain.Finding{}, e
	}
	rows, e := s.DB.SQL.Query(`SELECT id,inspection_id,title,severity,status,owner_id,description,due_at,resolution FROM findings WHERE id=?`, id)
	if e != nil {
		return domain.Finding{}, e
	}
	defer rows.Close()
	if !rows.Next() {
		return domain.Finding{}, errors.New("finding not found")
	}
	var f domain.Finding
	var due string
	e = rows.Scan(&f.ID, &f.InspectionID, &f.Title, &f.Severity, &f.Status, &f.OwnerID, &f.Description, &due, &f.Resolution)
	f.DueAt, _ = time.Parse(time.RFC3339, due)
	return f, e
}
