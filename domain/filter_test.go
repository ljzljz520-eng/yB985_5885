package domain

import (
	"testing"
	"time"
)

func TestFilter(t *testing.T) {
	i := Inspection{StoreID: "s", Status: StatusOpen, Summary: "leak", OpenedAt: time.Unix(1, 0), Findings: []Finding{{OwnerID: "u", Severity: "high"}}}
	if !(InspectionFilter{StoreID: "s", OwnerID: "u"}).Matches(i, Store{ID: "s"}) {
		t.Fatal()
	}
}
