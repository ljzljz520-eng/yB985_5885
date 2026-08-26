package query

import (
	"sort"
	"storeinspection/domain"
	"time"
)

type CalendarEntry struct {
	Date              time.Time
	Open, Due, Closed int
}

func Calendar(items []domain.Inspection) []CalendarEntry {
	m := map[string]*CalendarEntry{}
	for _, i := range items {
		k := i.OpenedAt.Format("2006-01-02")
		if m[k] == nil {
			m[k] = &CalendarEntry{Date: i.OpenedAt.Truncate(24 * time.Hour)}
		}
		m[k].Open++
		if i.Status == domain.StatusClosed {
			m[k].Closed++
		}
		for _, f := range i.Findings {
			if !f.DueAt.IsZero() {
				dk := f.DueAt.Format("2006-01-02")
				if m[dk] == nil {
					m[dk] = &CalendarEntry{Date: f.DueAt.Truncate(24 * time.Hour)}
				}
				m[dk].Due++
			}
		}
	}
	out := make([]CalendarEntry, 0, len(m))
	for _, e := range m {
		out = append(out, *e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Date.Before(out[j].Date) })
	return out
}
func InRange(entries []CalendarEntry, start, end time.Time) []CalendarEntry {
	var out []CalendarEntry
	for _, e := range entries {
		if !e.Date.Before(start) && e.Date.Before(end) {
			out = append(out, e)
		}
	}
	return out
}
func CalendarTotal(entries []CalendarEntry) (int, int, int) {
	o, d, c := 0, 0, 0
	for _, e := range entries {
		o += e.Open
		d += e.Due
		c += e.Closed
	}
	return o, d, c
}
