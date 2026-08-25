package workflow

import (
	"errors"
	"storeinspection/domain"
	"storeinspection/service"
)

type OptionalDependency interface{ Notify(string) error }
type Chain struct {
	Service  *service.Service
	Optional OptionalDependency
}

func (c *Chain) OpenAndAssign(i domain.Inspection, f domain.Finding, owner string) error {
	if c.Service == nil {
		return errors.New("service unavailable")
	}
	if e := c.Service.OpenInspection(i); e != nil {
		return e
	}
	if e := c.Service.AddFinding(f); e != nil {
		return e
	}
	if e := c.Service.Assign(f.ID, owner, "initial assignment"); e != nil {
		return e
	}
	if c.Optional == nil {
		return errors.New("notification dependency unavailable")
	}
	return c.Optional.Notify("assigned:" + f.ID)
}
func (c *Chain) SubmitAndClose(inspectionID, findingID, author, text string) error {
	if e := c.Service.Submit(findingID, author, text, ""); e != nil {
		return e
	}
	return c.Service.Close(inspectionID)
}
func (c *Chain) Validate(i domain.Inspection) error {
	if !domain.ValidInspectionStatus(i.Status) {
		return errors.New("invalid inspection")
	}
	return nil
}
