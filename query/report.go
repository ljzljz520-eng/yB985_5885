package query

import (
	"fmt"
	"storeinspection/domain"
)

func RenderText(r Result) string {
	out := "巡店记录\n"
	for _, i := range r.Items {
		out += fmt.Sprintf("%s %s %s\n", i.ID, i.StoreID, i.Status)
		for _, f := range i.Findings {
			out += fmt.Sprintf("  [%s] %s (%s)\n", f.Severity, f.Title, f.Status)
		}
	}
	return out
}
func FilterDescription(f domain.InspectionFilter) string {
	return fmt.Sprintf("store=%s status=%s region=%s owner=%s severity=%s", f.StoreID, f.Status, f.Region, f.OwnerID, f.Severity)
}
