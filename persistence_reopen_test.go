package main

import (
	"os"
	"storeinspection/domain"
	"storeinspection/store"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	p := "reopen.db"
	defer os.Remove(p)
	d, e := store.Open(p)
	if e != nil {
		t.Fatal(e)
	}
	if e = d.SaveStore(domain.Store{ID: "persist", Name: "P"}); e != nil {
		t.Fatal(e)
	}
	d.Close()
	d, e = store.Open(p)
	if e != nil {
		t.Fatal(e)
	}
	defer d.Close()
	if _, e = d.GetStore("persist"); e != nil {
		t.Fatal(e)
	}
}
