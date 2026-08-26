package query

import (
	"sort"
	"storeinspection/domain"
)

type Dashboard struct {
	Total                             int
	Open, Assigned, Submitted, Closed int
	Risk                              int
	ByStore                           map[string]int
}

func BuildDashboard(items []domain.Inspection) Dashboard {
	d := Dashboard{ByStore: map[string]int{}}
	for _, i := range items {
		d.Total++
		d.ByStore[i.StoreID]++
		switch i.Status {
		case domain.StatusOpen:
			d.Open++
		case domain.StatusAssigned:
			d.Assigned++
		case domain.StatusSubmitted:
			d.Submitted++
		case domain.StatusClosed:
			d.Closed++
		}
		d.Risk += domain.RiskScore(i.Findings)
	}
	return d
}
func TopStores(items []domain.Inspection, limit int) []string {
	counts := map[string]int{}
	for _, i := range items {
		counts[i.StoreID]++
	}
	type pair struct {
		s string
		n int
	}
	var ps []pair
	for s, n := range counts {
		ps = append(ps, pair{s, n})
	}
	sort.Slice(ps, func(i, j int) bool {
		if ps[i].n == ps[j].n {
			return ps[i].s < ps[j].s
		}
		return ps[i].n > ps[j].n
	})
	if limit > len(ps) {
		limit = len(ps)
	}
	out := make([]string, 0, limit)
	for _, p := range ps[:limit] {
		out = append(out, p.s)
	}
	return out
}
func Paginate(items []domain.Inspection, page, size int) []domain.Inspection {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	start := (page - 1) * size
	if start >= len(items) {
		return []domain.Inspection{}
	}
	end := start + size
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}
