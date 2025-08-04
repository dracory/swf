package swf

import (
	"fmt"
	"time"
)

// Example demonstrates how to use the workflow package
func Example() {
	// Create a new workflow
	wf := NewWorkflow()

	// Create steps
	step1 := NewStep("step1")
	step1.Title = "First Step"
	step1.Description = "This is the first step of the workflow"

	step2 := NewStep("step2")
	step2.Title = "Second Step"
	step2.Description = "This is the second step of the workflow"

	step3 := NewStep("step3")
	step3.Title = "Third Step"
	step3.Description = "This is the third step of the workflow"

	// Add steps to the workflow
	wf.AddStep(step1)
	wf.AddStep(step2)
	wf.AddStep(step3)

	// Get the current step
	currentStep := wf.GetCurrentStep()
	if currentStep != nil { // Added nil check for robustness, though unlikely here
		fmt.Printf("Current step: %s\n", currentStep.Name)
	} else {
		fmt.Println("Current step: <nil>")
	}

	// Check if a step is current
	isCurrent := wf.IsStepCurrent(step1)
	fmt.Printf("Is step1 current? %v\n", isCurrent)

	// Move to the next step
	err := wf.SetCurrentStep(step2)
	if err != nil {
		fmt.Printf("Error setting current step: %v\n", err)
	}

	// Check if a step is complete
	isComplete := wf.IsStepComplete(step1)
	fmt.Printf("Is step1 complete? %v\n", isComplete)

	// Get progress
	progress := wf.GetProgress()
	fmt.Printf("Progress: %d/%d steps completed (%.2f%%)\n",
		progress.Completed, progress.Total, progress.Percents)

	// Set metadata for a step
	wf.SetStepMeta(step2, "user", "john")
	wf.SetStepMeta(step2, "timestamp", time.Now().Format(time.RFC3339))

	// Get metadata for a step
	user := wf.GetStepMeta(step2, "user")
	fmt.Printf("User for step2: %v\n", user)

	// Mark a step as completed
	wf.MarkStepAsCompleted(step2)

	// Visualize the workflow
	dotGraph := wf.Visualize()
	fmt.Println("Workflow visualization (DOT format):")
	fmt.Println(dotGraph)

	// Serialize the workflow state
	state, err := wf.ToString()
	if err != nil {
		fmt.Printf("Error serializing workflow: %v\n", err)
	} else {
		fmt.Printf("Workflow state: %s\n", state)
	}

	// Create a new workflow and deserialize the state
	newWf := NewWorkflow()
	err = newWf.FromString(state)
	if err != nil {
		fmt.Printf("Error deserializing workflow: %v\n", err)
	} else {
		// --- FIX IS HERE ---
		// After FromString, newWf.steps is empty, so GetCurrentStep will return nil.
		// We need to check for nil before accessing fields.
		// We can still access the *name* of the current step from the restored state.
		fmt.Printf("Deserialized workflow current step name (from state): %s\n", newWf.GetState().CurrentStepName)
		deserializedCurrentStep := newWf.GetCurrentStep()
		if deserializedCurrentStep != nil {
			// This block will likely not be reached with the current FromString implementation
			fmt.Printf("Deserialized workflow current step object: %s\n", deserializedCurrentStep.Name)
		} else {
			fmt.Println("Deserialized workflow current step object: <nil> (Steps are not restored by FromString)")
		}
	}
}
