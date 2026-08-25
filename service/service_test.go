package service

import (
	"storeinspection/domain"
	"storeinspection/store"
	"testing"
	"time"
)

func TestServiceAssign(t *testing.T) {
	d, _ := store.Open(":memory:")
	defer d.Close()
	s := New(d, FixedClock{time.Unix(1, 0)})
	s.OpenInspection(domain.Inspection{ID: "i", StoreID: "s"})
	if e := s.AddFinding(domain.Finding{ID: "f", InspectionID: "i", Title: "x", Severity: "high"}); e != nil {
		t.Fatal(e)
	}
	if e := s.Assign("f", "u", ""); e != nil {
		t.Fatal(e)
	}
}
