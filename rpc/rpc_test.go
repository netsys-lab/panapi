package rpc

import (
	"testing"
)

func TestIDServer(t *testing.T) {
	s := IDServer{42}
	var id IDMsg
	s.GetID(new(IDMsg), &id)
	if *id.Value != 43 {
		t.Errorf("GetID Value = %d, want 43", *id.Value)
	}
	s.GetID(new(IDMsg), &id)
	if *id.Value != 44 {
		t.Errorf("GetID Value = %d, want 44", *id.Value)
	}

}
