package workflow

import (
	"storeinspection/domain"
	"storeinspection/service"
	"storeinspection/store"
	"testing"
)

func TestWorkflowOpenAssign(t *testing.T) {
	d, _ := store.Open(":memory:")
	defer d.Close()
	c := Chain{Service: service.New(d, service.FixedClock{}), Optional: &Recorder{}}
	if e := c.OpenAndAssign(domain.Inspection{ID: "i", StoreID: "s"}, domain.Finding{ID: "f", InspectionID: "i", Title: "x", Severity: "low"}, "u"); e != nil {
		t.Fatal(e)
	}
}
func TestWorkflowSubmitClose(t *testing.T) {
	d, _ := store.Open(":memory:")
	defer d.Close()
	s := service.New(d, service.FixedClock{})
	s.OpenInspection(domain.Inspection{ID: "i", StoreID: "s", Status: domain.StatusOpen})
	s.AddFinding(domain.Finding{ID: "f", InspectionID: "i", Title: "x", Severity: "low"})
	s.Assign("f", "u", "")
	c := Chain{Service: s, Optional: &Recorder{}}
	if e := c.SubmitAndClose("i", "f", "u", "fixed"); e != nil {
		t.Fatal(e)
	}
}
func TestWorkflowListFilter(t *testing.T) { d, _ := store.Open(":memory:"); defer d.Close() }
