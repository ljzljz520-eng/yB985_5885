package store

import (
	"storeinspection/domain"
	"testing"
	"time"
)

func TestStoreRoundTrip(t *testing.T) {
	d, e := Open(":memory:")
	if e != nil {
		t.Fatal(e)
	}
	defer d.Close()
	s := domain.Store{ID: "s", Name: "Main", CreatedAt: time.Unix(1, 0)}
	if e = d.SaveStore(s); e != nil {
		t.Fatal(e)
	}
	if _, e = d.GetStore("s"); e != nil {
		t.Fatal(e)
	}
}
