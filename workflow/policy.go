package workflow

import (
	"errors"
	"storeinspection/domain"
)

type Policy struct{ RequireOwner, RequireNote bool }

func (p Policy) CheckFinding(f domain.Finding) error {
	if p.RequireOwner && f.OwnerID == "" {
		return errors.New("finding owner required")
	}
	if p.RequireNote && f.Description == "" {
		return errors.New("finding description required")
	}
	return nil
}
func (p Policy) CanClose(fs []domain.Finding) bool {
	for _, f := range fs {
		if f.Status != domain.StatusSubmitted && f.Status != domain.StatusClosed {
			return false
		}
	}
	return true
}
func NormalizeOwner(owner string) string {
	if owner == "" {
		return "unassigned"
	}
	return owner
}
