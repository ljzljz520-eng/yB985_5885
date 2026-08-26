package domain

import "testing"

func TestValidation(t *testing.T) {
	if ValidateStore(Store{ID: "x", Name: "n"}) != nil {
		t.Fatal()
	}
	if !CanTransition(StatusOpen, StatusAssigned) {
		t.Fatal()
	}
}
