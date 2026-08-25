package domain

import "time"

type DailyMetric struct {
	Day                           time.Time
	Inspections, Findings, Closed int
}
type Trend struct {
	StoreID string
	Points  []DailyMetric
}

func BuildTrend(items []Inspection) Trend {
	out := Trend{}
	for _, i := range items {
		if out.StoreID == "" {
			out.StoreID = i.StoreID
		}
		day := i.OpenedAt.Truncate(24 * time.Hour)
		idx := -1
		for n := range out.Points {
			if out.Points[n].Day.Equal(day) {
				idx = n
			}
		}
		if idx < 0 {
			out.Points = append(out.Points, DailyMetric{Day: day})
			idx = len(out.Points) - 1
		}
		out.Points[idx].Inspections++
		out.Points[idx].Findings += len(i.Findings)
		if i.Status == StatusClosed {
			out.Points[idx].Closed++
		}
	}
	return out
}
func CompletionRate(m DailyMetric) float64 {
	if m.Inspections == 0 {
		return 0
	}
	return float64(m.Closed) / float64(m.Inspections)
}
func SeverityWeight(s string) int {
	switch s {
	case "high":
		return 3
	case "medium":
		return 2
	case "low":
		return 1
	default:
		return 0
	}
}
func RiskScore(fs []Finding) int {
	score := 0
	for _, f := range fs {
		if f.Status != StatusClosed {
			score += SeverityWeight(f.Severity)
		}
	}
	return score
}
