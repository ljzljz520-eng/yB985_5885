package workflow

import (
	"errors"
	"storeinspection/domain"
	"storeinspection/service"
)

type Review struct {
	Service  *service.Service
	Reviewer string
}

func (r Review) Approve(findingID string) error {
	if r.Service == nil {
		return errors.New("service unavailable")
	}
	f, e := r.ServiceLookup(findingID)
	if e != nil {
		return e
	}
	if f.Status != domain.StatusSubmitted {
		return errors.New("finding is not submitted")
	}
	if e = r.Service.DB.MarkApproved(findingID, true); e != nil {
		return e
	}
	return r.Service.DB.UpdateResolution(findingID, f.Resolution)
}
func (r Review) Reject(findingID, reason string) error {
	if reason == "" {
		return errors.New("rejection reason required")
	}
	f, e := r.ServiceLookup(findingID)
	if e != nil {
		return e
	}
	if f.Status != domain.StatusSubmitted {
		return errors.New("finding is not submitted")
	}
	f.Status = domain.StatusRejected
	f.Resolution = reason
	return r.Service.DB.SaveFinding(f)
}
func (r Review) ServiceLookup(id string) (domain.Finding, error) {
	rows, e := r.Service.DB.SQL.Query(`SELECT id,inspection_id,title,severity,status,owner_id,description,due_at,resolution FROM findings WHERE id=?`, id)
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
	return f, e
}
func (r Review) CanReview(f domain.Finding) bool {
	return f.Status == domain.StatusSubmitted && r.Reviewer != ""
}
