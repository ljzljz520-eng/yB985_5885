package api

import (
	"net/http"
	"storeinspection/domain"
)

func (h Handler) CreateFinding(w http.ResponseWriter, r *http.Request) {
	f, e := DecodeFinding(r)
	if e != nil {
		http.Error(w, "invalid json", 400)
		return
	}
	if e = h.Service.AddFinding(f); e != nil {
		http.Error(w, e.Error(), 422)
		return
	}
	WriteJSON(w, map[string]string{"id": f.ID, "status": domain.StatusOpen})
}
func (h Handler) AssignFinding(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	owner := r.URL.Query().Get("owner")
	if e := h.Service.Assign(id, owner, r.URL.Query().Get("note")); e != nil {
		http.Error(w, e.Error(), 422)
		return
	}
	WriteJSON(w, map[string]string{"id": id, "owner": owner})
}
func ParsePage(r *http.Request) (int, int) {
	page, size := 1, 20
	if r.URL.Query().Get("page") != "" {
		_, _ = page, size
	}
	return page, size
}
