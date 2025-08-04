package swf_test

import (
	"testing"

	"github.com/dracory/swf"
)

func TestNewStep(t *testing.T) {
	name := "test_step"
	step := swf.NewStep(name)

	if step == nil {
		t.Fatal("NewStep() returned nil")
	}

	if step.Name != name {
		t.Errorf("Expected name %s, got %s", name, step.Name)
	}

	if step.Type != "normal" {
		t.Errorf("Expected type 'normal', got %s", step.Type)
	}

	if step.Title != "" {
		t.Errorf("Expected empty title, got %s", step.Title)
	}

	if step.Description != "" {
		t.Errorf("Expected empty description, got %s", step.Description)
	}

	if step.Responsible != "Admin" {
		t.Errorf("Expected responsible 'Admin', got %s", step.Responsible)
	}
}

func TestGetActionLink(t *testing.T) {
	step := swf.NewStep("test_step")
	step.Responsible = "User"

	actionLink := step.GetActionLink()

	expectedLink := "/User/test_step"
	if actionLink != expectedLink {
		t.Errorf("Expected action link %s, got %s", expectedLink, actionLink)
	}

	// Test with different responsible
	step.Responsible = "Admin"
	actionLink = step.GetActionLink()

	expectedLink = "/Admin/test_step"
	if actionLink != expectedLink {
		t.Errorf("Expected action link %s, got %s", expectedLink, actionLink)
	}
}
