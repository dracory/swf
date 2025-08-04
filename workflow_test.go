package swf_test

import (
	"testing"

	"github.com/dracory/swf"
)

func TestNewWorkflow(t *testing.T) {
	wf := swf.NewWorkflow()

	if wf == nil {
		t.Fatal("NewWorkflow() returned nil")
	}

	if wf.GetState().CurrentStepName != "" {
		t.Errorf("Expected empty CurrentStepName, got %s", wf.GetState().CurrentStepName)
	}

	state := wf.GetState()
	if len(state.History) != 0 {
		t.Errorf("Expected empty History, got %v", state.History)
	}

	if state.StepDetails == nil {
		t.Error("StepDetails is nil")
	}

	steps := wf.GetSteps()
	if steps == nil {
		t.Error("steps is nil")
	}
}

func TestAddStep(t *testing.T) {
	wf := swf.NewWorkflow()
	step := swf.NewStep("test_step")

	err := wf.AddStep(step)
	if err != nil {
		t.Errorf("AddStep failed: %v", err)
	}

	steps := wf.GetSteps()
	if len(steps) != 1 {
		t.Errorf("Expected 1 step, got %d", len(steps))
	}

	if steps[0] != step {
		t.Error("Step not added correctly")
	}

	if wf.GetState().CurrentStepName != "test_step" {
		t.Errorf("Expected CurrentStepName to be 'test_step', got %s", wf.GetState().CurrentStepName)
	}

	if wf.GetState().StepDetails["test_step"].Completed != "" {
		t.Errorf("Expected Completed to be empty, got %s", wf.GetState().StepDetails["test_step"].Completed)
	}
}

func TestGetCurrentStep(t *testing.T) {
	wf := swf.NewWorkflow()

	// Test with no steps
	currentStep := wf.GetCurrentStep()
	if currentStep != nil {
		t.Error("Expected nil for current step with no steps")
	}

	// Test with steps
	step := swf.NewStep("test_step")
	err := wf.AddStep(step)
	if err != nil {
		t.Errorf("AddStep failed: %v", err)
	}

	currentStep = wf.GetCurrentStep()
	if currentStep == nil {
		t.Fatal("Expected non-nil current step")
	}

	if currentStep.Name != "test_step" {
		t.Errorf("Expected step name 'test_step', got %s", currentStep.Name)
	}
}

func TestSetCurrentStep(t *testing.T) {
	wf := swf.NewWorkflow()

	// Test with string
	step1 := swf.NewStep("step1")
	step2 := swf.NewStep("step2")
	err := wf.AddStep(step1)
	if err != nil {
		t.Errorf("AddStep failed: %v", err)
	}
	err = wf.AddStep(step2)
	if err != nil {
		t.Errorf("AddStep failed: %v", err)
	}

	err = wf.SetCurrentStep("step2")
	if err != nil {
		t.Errorf("SetCurrentStep failed: %v", err)
	}

	if wf.GetState().CurrentStepName != "step2" {
		t.Errorf("Expected CurrentStepName to be 'step2', got %s", wf.GetState().CurrentStepName)
	}

	// Test with Step pointer
	err = wf.SetCurrentStep(step1)
	if err != nil {
		t.Errorf("SetCurrentStep failed: %v", err)
	}

	if wf.GetState().CurrentStepName != "step1" {
		t.Errorf("Expected CurrentStepName to be 'step1', got %s", wf.GetState().CurrentStepName)
	}

	// Test with invalid step
	err = wf.SetCurrentStep("invalid_step")
	if err == nil {
		t.Error("Expected error for invalid step, got nil")
	}

	// Test with invalid type
	err = wf.SetCurrentStep(123)
	if err == nil {
		t.Error("Expected error for invalid type, got nil")
	}
}

func TestIsStepCurrent(t *testing.T) {
	wf := swf.NewWorkflow()
	step1 := swf.NewStep("step1")
	step2 := swf.NewStep("step2")
	err := wf.AddStep(step1)
	if err != nil {
		t.Errorf("AddStep failed: %v", err)
	}
	err = wf.AddStep(step2)
	if err != nil {
		t.Errorf("AddStep failed: %v", err)
	}

	// Test with string
	isCurrent := wf.IsStepCurrent("step1")
	if !isCurrent {
		t.Error("Expected step1 to be current")
	}

	// Test with Step pointer
	isCurrent = wf.IsStepCurrent(step2)
	if isCurrent {
		t.Error("Expected step2 to not be current")
	}

	// Test with invalid type
	isCurrent = wf.IsStepCurrent(123)
	if isCurrent {
		t.Error("Expected false for invalid type, got true")
	}
}

func TestIsStepComplete(t *testing.T) {
	wf := swf.NewWorkflow()
	step1 := swf.NewStep("step1")
	step2 := swf.NewStep("step2")
	step3 := swf.NewStep("step3")
	err := wf.AddStep(step1)
	if err != nil {
		t.Errorf("AddStep failed: %v", err)
	}
	err = wf.AddStep(step2)
	if err != nil {
		t.Errorf("AddStep failed: %v", err)
	}
	err = wf.AddStep(step3)
	if err != nil {
		t.Errorf("AddStep failed: %v", err)
	}

	// Initially, no steps should be complete
	if wf.IsStepComplete(step1) {
		t.Error("Expected step1 to not be complete initially")
	}

	// Move to step2
	err = wf.SetCurrentStep(step2)
	if err != nil {
		t.Errorf("SetCurrentStep failed: %v", err)
	}

	// Now step1 should be complete
	if !wf.IsStepComplete(step1) {
		t.Error("Expected step1 to be complete after moving to step2")
	}

	// step2 should not be complete
	if wf.IsStepComplete(step2) {
		t.Error("Expected step2 to not be complete")
	}

	// Mark step2 as completed
	isMarked := wf.MarkStepAsCompleted(step2)
	if !isMarked {
		t.Error("Expected step2 to be marked as completed")
	}

	// Now step2 should be complete
	if !wf.IsStepComplete(step2) {
		t.Error("Expected step2 to be complete after marking as completed")
	}

	// Test with invalid type
	if wf.IsStepComplete(123) {
		t.Error("Expected false for invalid type, got true")
	}
}

func TestGetProgress(t *testing.T) {
	wf := swf.NewWorkflow()
	step1 := swf.NewStep("step1")
	step2 := swf.NewStep("step2")
	step3 := swf.NewStep("step3")
	err := wf.AddStep(step1)
	if err != nil {
		t.Errorf("AddStep failed: %v", err)
	}
	err = wf.AddStep(step2)
	if err != nil {
		t.Errorf("AddStep failed: %v", err)
	}
	err = wf.AddStep(step3)
	if err != nil {
		t.Errorf("AddStep failed: %v", err)
	}

	// Initially, 0 steps should be complete
	progress := wf.GetProgress()
	if progress.Total != 3 {
		t.Errorf("Expected 3 total steps, got %d", progress.Total)
	}
	if progress.Completed != 0 {
		t.Errorf("Expected 0 completed steps, got %d", progress.Completed)
	}
	if progress.Current != 0 {
		t.Errorf("Expected current step 0, got %d", progress.Current)
	}
	if progress.Pending != 3 {
		t.Errorf("Expected 3 pending steps, got %d", progress.Pending)
	}
	if progress.Percents != 0 {
		t.Errorf("Expected 0%% progress, got %.2f%%", progress.Percents)
	}

	// Move to step2
	err = wf.SetCurrentStep(step2)
	if err != nil {
		t.Errorf("SetCurrentStep failed: %v", err)
	}

	// Now 1 step should be complete
	progress = wf.GetProgress()
	if progress.Completed != 1 {
		t.Errorf("Expected 1 completed step, got %d", progress.Completed)
	}
	if progress.Current != 1 {
		t.Errorf("Expected current step 1, got %d", progress.Current)
	}
	if progress.Pending != 2 {
		t.Errorf("Expected 2 pending steps, got %d", progress.Pending)
	}
	if progress.Percents != 33.33333333333333 {
		t.Errorf("Expected 33.33%% progress, got %.2f%%", progress.Percents)
	}

	// Mark step2 as completed
	isMarked := wf.MarkStepAsCompleted(step2)
	if !isMarked {
		t.Error("Expected step2 to be marked as completed")
	}

	// Now 2 steps should be complete
	progress = wf.GetProgress()
	if progress.Completed != 2 {
		t.Errorf("Expected 2 completed steps, got %d", progress.Completed)
	}
	if progress.Pending != 1 {
		t.Errorf("Expected 1 pending step, got %d", progress.Pending)
	}
	if progress.Percents != 66.66666666666666 {
		t.Errorf("Expected 66.67%% progress, got %.2f%%", progress.Percents)
	}
}

func TestGetSteps(t *testing.T) {
	wf := swf.NewWorkflow()
	step1 := swf.NewStep("step1")
	step2 := swf.NewStep("step2")
	err := wf.AddStep(step1)
	if err != nil {
		t.Errorf("AddStep failed: %v", err)
	}
	err = wf.AddStep(step2)
	if err != nil {
		t.Errorf("AddStep failed: %v", err)
	}

	steps := wf.GetSteps()

	if len(steps) != 2 {
		t.Errorf("Expected 2 steps, got %d", len(steps))
	}

	if steps[0] != step1 {
		t.Error("step1 not found in steps")
	}

	if steps[1] != step2 {
		t.Error("step2 not found in steps")
	}
}

func TestGetStep(t *testing.T) {
	wf := swf.NewWorkflow()
	step1 := swf.NewStep("step1")
	err := wf.AddStep(step1)
	if err != nil {
		t.Errorf("AddStep failed: %v", err)
	}

	// Test with existing step
	step := wf.GetStep("step1")
	if step != step1 {
		t.Error("GetStep returned incorrect step")
	}

	// Test with non-existing step
	step = wf.GetStep("non_existing")
	if step != nil {
		t.Error("Expected nil for non-existing step, got non-nil")
	}
}

func TestGetStepMeta(t *testing.T) {
	wf := swf.NewWorkflow()
	step := swf.NewStep("test_step")
	err := wf.AddStep(step)
	if err != nil {
		t.Errorf("AddStep failed: %v", err)
	}

	// Test with no metadata
	meta := wf.GetStepMeta(step, "key")
	if meta != nil {
		t.Error("Expected nil for non-existing metadata, got non-nil")
	}

	// Set metadata
	wf.SetStepMeta(step, "key", "value")

	// Test with existing metadata
	meta = wf.GetStepMeta(step, "key")
	if meta != "value" {
		t.Errorf("Expected 'value' for metadata, got %v", meta)
	}

	// Test with string step name
	meta = wf.GetStepMeta("test_step", "key")
	if meta != "value" {
		t.Errorf("Expected 'value' for metadata with string step name, got %v", meta)
	}

	// Test with invalid type
	meta = wf.GetStepMeta(123, "key")
	if meta != nil {
		t.Error("Expected nil for invalid step type, got non-nil")
	}
}

func TestMarkStepAsCompleted(t *testing.T) {
	wf := swf.NewWorkflow()
	step := swf.NewStep("test_step")
	err := wf.AddStep(step)
	if err != nil {
		t.Errorf("AddStep failed: %v", err)
	}

	// Test with non-existing step details
	result := wf.MarkStepAsCompleted("non_existing")
	if result {
		t.Error("Expected false for non-existing step, got true")
	}

	// Test with existing step
	result = wf.MarkStepAsCompleted(step)
	if !result {
		t.Error("Expected true for existing step, got false")
	}

	// Check if step is marked as completed
	state := wf.GetState()
	details := state.StepDetails[step.Name]
	if details.Completed == "" {
		t.Error("Step not marked as completed")
	}

	// Test with string step name
	result = wf.MarkStepAsCompleted("test_step")
	if !result {
		t.Error("Expected true for string step name, got false")
	}

	// Test with invalid type
	result = wf.MarkStepAsCompleted(123)
	if result {
		t.Error("Expected false for invalid step type, got true")
	}
}

func TestSetStepMeta(t *testing.T) {
	wf := swf.NewWorkflow()
	step := swf.NewStep("test_step")
	err := wf.AddStep(step)
	if err != nil {
		t.Errorf("AddStep failed: %v", err)
	}

	// Set metadata
	wf.SetStepMeta(step, "key", "value")

	// Check if metadata is set
	state := wf.GetState()
	details := state.StepDetails[step.Name]
	if details.Meta == nil {
		t.Fatal("Meta map is nil")
	}

	value, exists := details.Meta["key"]
	if !exists {
		t.Error("Metadata not set")
	}

	if value != "value" {
		t.Errorf("Expected 'value' for metadata, got %v", value)
	}

	// Test with string step name
	wf.SetStepMeta("test_step", "key2", "value2")

	state = wf.GetState()
	details = state.StepDetails[step.Name]
	value, exists = details.Meta["key2"]
	if !exists {
		t.Error("Metadata not set with string step name")
	}

	if value != "value2" {
		t.Errorf("Expected 'value2' for metadata with string step name, got %v", value)
	}

	// Test with invalid type
	wf.SetStepMeta(123, "key", "value")
	// Should not panic
}

func TestGetState(t *testing.T) {
	wf := swf.NewWorkflow()
	step := swf.NewStep("test_step")
	err := wf.AddStep(step)
	if err != nil {
		t.Errorf("AddStep failed: %v", err)
	}
	step2 := swf.NewStep("test_step2")
	err = wf.AddStep(step2)
	if err != nil {
		t.Errorf("AddStep failed: %v", err)
	}

	state := wf.GetState()

	if state.CurrentStepName != "test_step" {
		t.Errorf("Expected CurrentStepName to be 'test_step', got %s", state.CurrentStepName)
	}

	if state.StepDetails["test_step"].Started == "" {
		t.Errorf("Expected Started to not be empty, got %s", state.StepDetails["test_step"].Started)
	}

	if state.StepDetails["test_step"].Completed != "" {
		t.Errorf("Expected Completed to be empty, got %s", state.StepDetails["test_step"].Completed)
	}

	if len(state.History) != 1 {
		t.Errorf("Expected 1 history entry (current step, which has been started), got %d", len(state.History))
	}

	isMarked := wf.MarkStepAsCompleted(step)
	if !isMarked {
		t.Error("Expected step to be marked as completed")
	}

	state = wf.GetState()

	if len(state.History) != 1 {
		t.Errorf("Expected 1 history entry (current step, which has been completed), got %d", len(state.History))
	}

	if state.History[0] != "test_step" {
		t.Errorf("Expected history entry to be 'test_step', got %s", state.History[0])
	}

	if state.StepDetails["test_step"].Completed == "" {
		t.Errorf("Expected Completed to not be empty, got %s", state.StepDetails["test_step"].Completed)
	}
}

func TestToStringAndFromString(t *testing.T) {
	wf := swf.NewWorkflow()
	step1 := swf.NewStep("step1")
	step2 := swf.NewStep("step2")
	err := wf.AddStep(step1)
	if err != nil {
		t.Errorf("AddStep failed: %v", err)
	}
	err = wf.AddStep(step2)
	if err != nil {
		t.Errorf("AddStep failed: %v", err)
	}

	// Set some metadata
	wf.SetStepMeta(step1, "key", "value")

	// Serialize
	str, err := wf.ToString()
	if err != nil {
		t.Fatalf("ToString failed: %v", err)
	}

	// Create a new workflow and deserialize
	newWf := swf.NewWorkflow()
	err = newWf.FromString(str)
	if err != nil {
		t.Fatalf("FromString failed: %v", err)
	}

	// Check if state is the same
	if newWf.GetState().CurrentStepName != wf.GetState().CurrentStepName {
		t.Errorf("Expected CurrentStepName %s, got %s", wf.GetState().CurrentStepName, newWf.GetState().CurrentStepName)
	}

	if len(newWf.GetState().History) != len(wf.GetState().History) {
		t.Errorf("Expected %d history entries, got %d", len(wf.GetState().History), len(newWf.GetState().History))
	}

	// Check metadata
	meta := newWf.GetStepMeta(step1, "key")
	if meta != "value" {
		t.Errorf("Expected metadata 'value', got %v", meta)
	}

	// Test with invalid JSON
	err = newWf.FromString("invalid json")
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}
