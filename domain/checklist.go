package domain

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

type ChecklistItem struct {
	Code, Label, Guidance string
	Required              bool
	Weight                int
}
type Checklist struct {
	ID, Name, Version string
	Items             []ChecklistItem
}
type CheckResult struct {
	ItemCode, Value, Note string
	Score                 int
	At                    time.Time
}

func DefaultChecklist() Checklist {
	return Checklist{ID: "retail-standard", Name: "门店标准巡检", Version: "2026.1", Items: []ChecklistItem{{"safety", "消防通道", "保持畅通", true, 5}, {"display", "陈列规范", "按图陈列", true, 3}, {"stock", "库存准确", "盘点差异", false, 2}, {"service", "服务礼仪", "主动问候", true, 4}}}
}
func (c Checklist) Validate() error {
	if c.ID == "" || c.Name == "" {
		return errors.New("checklist identity required")
	}
	seen := map[string]bool{}
	for _, i := range c.Items {
		if i.Code == "" || i.Label == "" {
			return errors.New("checklist item incomplete")
		}
		if seen[i.Code] {
			return fmt.Errorf("duplicate item %s", i.Code)
		}
		seen[i.Code] = true
		if i.Weight < 0 {
			return errors.New("negative weight")
		}
	}
	return nil
}
func (c Checklist) RequiredCodes() []string {
	var out []string
	for _, i := range c.Items {
		if i.Required {
			out = append(out, i.Code)
		}
	}
	sort.Strings(out)
	return out
}
func (c Checklist) Score(results []CheckResult) int {
	weights := map[string]int{}
	for _, i := range c.Items {
		weights[i.Code] = i.Weight
	}
	total, possible := 0, 0
	for _, r := range results {
		if w, ok := weights[r.ItemCode]; ok {
			total += r.Score * w
			possible += w * 5
		}
	}
	if possible == 0 {
		return 0
	}
	return total * 100 / possible
}
func CleanNote(s string) string { return strings.Join(strings.Fields(s), " ") }
