package query

import (
	"storeinspection/domain"
	"storeinspection/store"
	"testing"
)

func TestSearch(t *testing.T) {
	d, _ := store.Open(":memory:")
	defer d.Close()
	d.SaveStore(domain.Store{ID: "s", Name: "A"})
	d.SaveInspection(domain.Inspection{ID: "i", StoreID: "s", Status: domain.StatusOpen})
	r, e := Search(d, domain.InspectionFilter{Status: "open"})
	if e != nil || len(r.Items) != 1 {
		t.Fatal(e)
	}
}
