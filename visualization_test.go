package swf

import (
	"strings"
	"testing"
)

func TestVisualize(t *testing.T) {
	// Create a workflow
	wf := NewWorkflow()

	// Create steps
	step1 := NewStep("step1")
	step1.Title = "Document Review"
	step1.Description = "Review the submitted document"

	step2 := NewStep("step2")
	step2.Title = "Manager Approval"
	step2.Description = "Manager approval of the document"

	step3 := NewStep("step3")
	step3.Title = "Final Sign-off"
	step3.Description = "Final sign-off by executive"

	// Add steps
	wf.AddStep(step1)
	wf.AddStep(step2)
	wf.AddStep(step3)

	// Generate DOT graph
	dot := wf.Visualize()

	// Basic validation
	if !strings.Contains(dot, "digraph") {
		t.Error("Expected DOT graph to contain 'digraph'")
	}

	// Check if all steps are present
	if !strings.Contains(dot, "step1") || !strings.Contains(dot, "step2") || !strings.Contains(dot, "step3") {
		t.Error("Expected DOT graph to contain all steps")
	}

	// Check if edges are present
	if !strings.Contains(dot, "->") {
		t.Error("Expected DOT graph to contain edges")
	}

	// Move to step2 and check visualization
	wf.SetCurrentStep(step2)
	dot = wf.Visualize()

	// Step1 should be green (completed)
	if !strings.Contains(dot, `fillcolor="#4CAF50"`) {
		t.Error("Expected completed step to have green fill color")
	}

	// Step2 should be blue (current)
	if !strings.Contains(dot, `fillcolor="#2196F3"`) {
		t.Error("Expected current step to have blue fill color")
	}

	t.Log(dot)

	// Test with empty workflow
	emptyWf := NewWorkflow()
	emptyDot := emptyWf.Visualize()
	if !strings.Contains(emptyDot, "digraph") {
		t.Error("Expected empty workflow to generate valid DOT graph")
	}
}
