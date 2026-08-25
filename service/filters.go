package service

import (
	"errors"
	"sort"
	"storeinspection/domain"
	"storeinspection/store"
	"strings"
)

type FilterPreset struct {
	ID, Name string
	Filter   domain.InspectionFilter
	Owner    string
}

func SavePreset(presets map[string]FilterPreset, p FilterPreset) error {
	if strings.TrimSpace(p.ID) == "" || strings.TrimSpace(p.Name) == "" {
		return errors.New("preset identity required")
	}
	presets[p.ID] = p
	return nil
}
func DeletePreset(presets map[string]FilterPreset, id string) bool {
	if _, ok := presets[id]; !ok {
		return false
	}
	delete(presets, id)
	return true
}
func ListPresets(presets map[string]FilterPreset) []FilterPreset {
	out := make([]FilterPreset, 0, len(presets))
	for _, p := range presets {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
func ApplyPreset(db *store.DB, p FilterPreset) ([]domain.Inspection, error) {
	items, e := db.ListInspections()
	if e != nil {
		return nil, e
	}
	stores, e := db.ListStores()
	if e != nil {
		return nil, e
	}
	by := map[string]domain.Store{}
	for _, s := range stores {
		by[s.ID] = s
	}
	var out []domain.Inspection
	for _, i := range items {
		if p.Filter.Matches(i, by[i.StoreID]) {
			out = append(out, i)
		}
	}
	return out, nil
}
func NormalizePresetName(n string) string { return strings.Join(strings.Fields(n), " ") }
