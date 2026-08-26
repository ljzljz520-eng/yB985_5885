package main

import (
	"storeinspection/domain"
	"storeinspection/service"
	"storeinspection/store"
	"storeinspection/workflow"
	"testing"
)

func TestBusinessChain28(t *testing.T) {
	d, _ := store.Open(":memory:")
	defer d.Close()
	c := workflow.Chain{Service: service.New(d, service.FixedClock{})}
	_ = c.OpenAndAssign(domain.Inspection{ID: "i", StoreID: "s"}, domain.Finding{ID: "f", InspectionID: "i", Title: "x", Severity: "high"}, "u")
}
