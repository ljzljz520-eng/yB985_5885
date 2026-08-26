package query

import (
	"storeinspection/domain"
	"storeinspection/store"
)

type Result struct {
	Items  []domain.Inspection
	Counts map[string]int
}

func Search(db *store.DB, f domain.InspectionFilter) (Result, error) {
	items, e := db.ListInspections()
	if e != nil {
		return Result{}, e
	}
	stores, e := db.ListStores()
	if e != nil {
		return Result{}, e
	}
	by := map[string]domain.Store{}
	for _, s := range stores {
		by[s.ID] = s
	}
	out := Result{Counts: map[string]int{}}
	for _, i := range items {
		if f.Matches(i, by[i.StoreID]) {
			fs, _ := db.ListFindings(i.ID)
			i.Findings = fs
			out.Items = append(out.Items, i)
			out.Counts[i.Status]++
		}
	}
	domain.SortInspections(out.Items, true)
	return out, nil
}
func StatusCounts(items []domain.Inspection) map[string]int {
	m := map[string]int{}
	for _, i := range items {
		m[i.Status]++
	}
	return m
}
