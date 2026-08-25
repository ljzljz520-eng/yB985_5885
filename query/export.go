package query

import (
	"encoding/json"
	"sort"
	"storeinspection/domain"
)

type ExportRow struct{ InspectionID, StoreID, Status, FindingID, Severity, Owner, Resolution string }

func Flatten(items []domain.Inspection) []ExportRow {
	var out []ExportRow
	for _, i := range items {
		if len(i.Findings) == 0 {
			out = append(out, ExportRow{InspectionID: i.ID, StoreID: i.StoreID, Status: i.Status})
			continue
		}
		for _, f := range i.Findings {
			out = append(out, ExportRow{i.ID, i.StoreID, i.Status, f.ID, f.Severity, f.OwnerID, f.Resolution})
		}
	}
	sort.Slice(out, func(a, b int) bool { return out[a].InspectionID < out[b].InspectionID })
	return out
}
func JSON(items []domain.Inspection) (string, error) {
	b, e := json.Marshal(Flatten(items))
	return string(b), e
}
func GroupByStore(items []domain.Inspection) map[string][]domain.Inspection {
	m := map[string][]domain.Inspection{}
	for _, i := range items {
		m[i.StoreID] = append(m[i.StoreID], i)
	}
	return m
}
