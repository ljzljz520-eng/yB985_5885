package query

import (
	"sort"
	"storeinspection/domain"
	"strings"
)

type Options struct{ Statuses, Regions, Severities, Owners []string }

func BuildOptions(items []domain.Inspection) Options {
	set := func() map[string]bool { return map[string]bool{} }
	ss, _, sv, oo := set(), set(), set(), set()
	for _, i := range items {
		ss[i.Status] = true
		for _, f := range i.Findings {
			sv[f.Severity] = true
			if f.OwnerID != "" {
				oo[f.OwnerID] = true
			}
		}
	}
	out := Options{}
	for k := range ss {
		out.Statuses = append(out.Statuses, k)
	}
	for k := range sv {
		out.Severities = append(out.Severities, k)
	}
	for k := range oo {
		out.Owners = append(out.Owners, k)
	}
	sort.Strings(out.Statuses)
	sort.Strings(out.Severities)
	sort.Strings(out.Owners)
	return out
}
func MatchAny(value string, choices []string) bool {
	if len(choices) == 0 {
		return true
	}
	for _, c := range choices {
		if strings.EqualFold(value, c) {
			return true
		}
	}
	return false
}
func Apply(items []domain.Inspection, statuses []string) []domain.Inspection {
	var out []domain.Inspection
	for _, i := range items {
		if MatchAny(i.Status, statuses) {
			out = append(out, i)
		}
	}
	return out
}
func DistinctStores(items []domain.Inspection) []string {
	m := map[string]bool{}
	for _, i := range items {
		m[i.StoreID] = true
	}
	var out []string
	for s := range m {
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}
