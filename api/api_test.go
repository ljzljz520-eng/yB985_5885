package api

import (
	"net/http/httptest"
	"storeinspection/service"
	"storeinspection/store"
	"testing"
)

func TestHealth(t *testing.T) {
	d, _ := store.Open(":memory:")
	defer d.Close()
	r := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	(Handler{Service: service.New(d, service.FixedClock{})}).ServeHTTP(r, req)
	if r.Code != 200 {
		t.Fatal(r.Code)
	}
}
