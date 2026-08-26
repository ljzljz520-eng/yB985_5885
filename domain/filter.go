package domain

import "strings"

type InspectionFilter struct{ StoreID, Status, Region, OwnerID, Severity, Query string }

func (f InspectionFilter) Matches(i Inspection, s Store) bool {
	if f.StoreID != "" && i.StoreID != f.StoreID {
		return false
	}
	if f.Status != "" && NormalizeStatus(i.Status) != NormalizeStatus(f.Status) {
		return false
	}
	if f.Region != "" && s.Region != f.Region {
		return false
	}
	if f.Query != "" {
		q := strings.ToLower(f.Query)
		if !strings.Contains(strings.ToLower(i.Summary), q) && !strings.Contains(strings.ToLower(i.Inspector), q) {
			return false
		}
	}
	if f.OwnerID != "" {
		found := false
		for _, x := range i.Findings {
			if x.OwnerID == f.OwnerID {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	if f.Severity != "" {
		found := false
		for _, x := range i.Findings {
			if x.Severity == f.Severity {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}
func SortInspections(items []Inspection, descending bool) {
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			swap := items[j].OpenedAt.Before(items[i].OpenedAt)
			if descending {
				swap = !swap
			}
			if swap {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}
