package api

import (
	"net/http"
	"strconv"
	"strings"
)

func Header(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
}
func ParseBool(v string, defaultValue bool) bool {
	v = strings.ToLower(strings.TrimSpace(v))
	if v == "" {
		return defaultValue
	}
	return v == "1" || v == "true" || v == "yes"
}
func ParseInt(v string, defaultValue, max int) int {
	n, e := strconv.Atoi(v)
	if e != nil || n < 0 {
		return defaultValue
	}
	if max > 0 && n > max {
		return max
	}
	return n
}
func CSV(values []string) string { return strings.Join(values, ",") }
func SplitCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
