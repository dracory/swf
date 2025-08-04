package swf_test

import (
	"testing"

	"github.com/dracory/swf"
)

func TestNewState(t *testing.T) {
	name := "test_state"
	state := swf.NewState(name)

	if state == nil {
		t.Fatal("NewState() returned nil")
	}

	if state.Name != name {
		t.Errorf("Expected name %s, got %s", name, state.Name)
	}

	if state.Title != "" {
		t.Errorf("Expected empty title, got %s", state.Title)
	}

	if state.Description != "" {
		t.Errorf("Expected empty description, got %s", state.Description)
	}

	if state.Responsible != "Admin" {
		t.Errorf("Expected responsible 'Admin', got %s", state.Responsible)
	}
}

func TestStateGetActionLink(t *testing.T) {
	state := swf.NewState("test_state")
	state.Responsible = "User"

	actionLink := state.GetActionLink()

	expectedLink := "/User/test_state"
	if actionLink != expectedLink {
		t.Errorf("Expected action link %s, got %s", expectedLink, actionLink)
	}

	// Test with different responsible
	state.Responsible = "Admin"
	actionLink = state.GetActionLink()

	expectedLink = "/Admin/test_state"
	if actionLink != expectedLink {
		t.Errorf("Expected action link %s, got %s", expectedLink, actionLink)
	}
}
