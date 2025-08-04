package swf_test

import (
	"strings"
	"testing"

	"github.com/dracory/swf"
)

func TestExample(t *testing.T) {
	// This test simply ensures that the Example function runs without panicking
	// It doesn't check the output since it's just printing to stdout
	swf.Example()
}

func TestExampleOutput(t *testing.T) {
	// Create a new workflow
	wf := swf.NewWorkflow()

	// Create steps
	step1 := swf.NewStep("step1")
	step1.Title = "First Step"
	step1.Description = "This is the first step of the workflow"

	step2 := swf.NewStep("step2")
	step2.Title = "Second Step"
	step2.Description = "This is the second step of the workflow"

	step3 := swf.NewStep("step3")
	step3.Title = "Third Step"
	step3.Description = "This is the third step of the workflow"

	// Add steps to the workflow
	wf.AddStep(step1)
	wf.AddStep(step2)
	wf.AddStep(step3)

	// Get the current step
	currentStep := wf.GetCurrentStep()
	if currentStep == nil {
		t.Fatal("Expected non-nil current step")
	}
	if currentStep.Name != "step1" {
		t.Errorf("Expected current step name 'step1', got %s", currentStep.Name)
	}

	// Check if a step is current
	isCurrent := wf.IsStepCurrent(step1)
	if !isCurrent {
		t.Error("Expected step1 to be current")
	}

	// Move to the next step
	err := wf.SetCurrentStep(step2)
	if err != nil {
		t.Errorf("Error setting current step: %v", err)
	}

	// Check if a step is complete
	isComplete := wf.IsStepComplete(step1)
	if !isComplete {
		t.Error("Expected step1 to be complete")
	}

	// Get progress
	progress := wf.GetProgress()
	if progress.Total != 3 {
		t.Errorf("Expected 3 total steps, got %d", progress.Total)
	}
	if progress.Completed != 1 {
		t.Errorf("Expected 1 completed step, got %d", progress.Completed)
	}

	// Set metadata for a step
	wf.SetStepMeta(step2, "user", "john")

	// Get metadata for a step
	user := wf.GetStepMeta(step2, "user")
	if user != "john" {
		t.Errorf("Expected user 'john', got %v", user)
	}

	// Mark a step as completed
	wf.MarkStepAsCompleted(step2)

	// Serialize the workflow state
	state, err := wf.ToString()
	if err != nil {
		t.Errorf("Error serializing workflow: %v", err)
	}
	if !strings.Contains(state, "step2") {
		t.Error("Expected state to contain 'step2'")
	}

	// Create a new workflow and deserialize the state
	newWf := swf.NewWorkflow()
	err = newWf.FromString(state)
	if err != nil {
		t.Errorf("Error deserializing workflow: %v", err)
	}

	// Note: After deserialization, the steps map is empty, so GetCurrentStep will return nil
	// This is expected behavior since we only serialize the state, not the steps
}
