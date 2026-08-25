package api

import (
	"encoding/json"
	"net/http"
	"storeinspection/domain"
	"storeinspection/query"
	"storeinspection/service"
)

type Handler struct{ Service *service.Service }

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/health":
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	case "/inspections":
		h.list(w, r)
	default:
		http.NotFound(w, r)
	}
}
func (h Handler) list(w http.ResponseWriter, r *http.Request) {
	f := domain.InspectionFilter{StoreID: r.URL.Query().Get("store"), Status: r.URL.Query().Get("status"), Region: r.URL.Query().Get("region"), OwnerID: r.URL.Query().Get("owner"), Severity: r.URL.Query().Get("severity"), Query: r.URL.Query().Get("q")}
	res, e := query.Search(h.Service.DB, f)
	if e != nil {
		http.Error(w, e.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(res)
}
func DecodeFinding(r *http.Request) (domain.Finding, error) {
	var f domain.Finding
	e := json.NewDecoder(r.Body).Decode(&f)
	return f, e
}
func WriteJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
