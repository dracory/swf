package swf

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dracory/base/arr"
	"github.com/samber/lo"
)

// StepDetails contains metadata about a step
type StepDetails struct {
	Started   string
	Completed string
	Meta      map[string]any
}

// WorkflowState represents the current state of a workflow
type WorkflowState struct {
	CurrentStepName string
	// History is the history of steps that have been completed
	// and the current step, which has been started
	History     []string
	StepDetails map[string]*StepDetails
}

// Progress represents workflow progress
type Progress struct {
	Total     int
	Completed int
	Current   int
	Pending   int
	Percents  float64
}

// Workflow represents a workflow
type Workflow struct {
	steps []*Step
	state *WorkflowState
}

// NewWorkflow creates a new Workflow
func NewWorkflow() *Workflow {
	return &Workflow{
		steps: make([]*Step, 0),
		state: &WorkflowState{
			History:     make([]string, 0),
			StepDetails: make(map[string]*StepDetails),
		},
	}
}

// AddStep adds a step to the workflow
//
// Business logic:
// 1. Check if step already exists
// 2. Add step to steps map
// 3. If first step, set it as current step
// 4. Add step details to step details map
func (w *Workflow) AddStep(step *Step) error {
	if w.GetStep(step.Name) != nil {
		return fmt.Errorf("step already exists: %s", step.Name)
	}

	w.steps = append(w.steps, step)

	w.state.StepDetails[step.Name] = &StepDetails{
		Started:   "",
		Completed: "",
		Meta:      make(map[string]any),
	}

	// if first step becomes current step
	if w.state.CurrentStepName == "" {
		w.SetCurrentStep(step.Name)
	}

	return nil
}

// GetCurrentStep returns the current step
func (w *Workflow) GetCurrentStep() *Step {
	if w.state.CurrentStepName == "" {
		return nil
	}

	return w.GetStep(w.state.CurrentStepName)
}

// SetCurrentStep sets the current step, can be a step name or a step pointer
//
// Business logic:
// 1. Check if step exists
// 2. Mark the current step as completed
// 3. Set the current step to the new step
func (w *Workflow) SetCurrentStep(step any) error {
	stepName, err := stepName(step)
	if err != nil {
		return err
	}

	if w.GetStep(stepName) == nil {
		return fmt.Errorf("step not found: %s", stepName)
	}

	// Mark the current step as completed
	if w.state.CurrentStepName != "" && w.state.CurrentStepName != stepName {
		w.state.StepDetails[w.state.CurrentStepName].Completed = time.Now().Format(time.RFC3339)
	}

	w.state.CurrentStepName = stepName
	w.state.History = append(w.state.History, stepName)
	w.state.StepDetails[stepName].Started = time.Now().Format(time.RFC3339)
	return nil
}

// IsStepCurrent checks if a step is the current step
func (w *Workflow) IsStepCurrent(step any) bool {
	stepName, err := stepName(step)
	if err != nil {
		return false
	}

	return w.state.CurrentStepName == stepName
}

// IsStepComplete checks if a step is completed
//
// Business logic:
// 1. Get step name
// 2. Get step positions
// 3. If step is before the current step, it's complete
// 4. If step is explicitly marked as completed, it's complete
func (w *Workflow) IsStepComplete(step any) bool {
	stepName, err := stepName(step)
	if err != nil {
		return false
	}

	// Get step positions
	stepNames := lo.Map(w.steps, func(item *Step, index int) string {
		return item.Name
	})

	currentStepPosition := arr.Index(stepNames, w.state.CurrentStepName)
	stepPosition := arr.Index(stepNames, stepName)

	// If the step is before the current step, it's complete
	if stepPosition < currentStepPosition {
		return true
	}

	// Check if the step is explicitly marked as completed
	return w.state.StepDetails[stepName].Completed != ""
}

// stepName returns the name of a step, can be a step name or a step pointer
//
// Business logic:
// 1. Check if step is a string
// 2. Check if step is a step pointer
// 3. Return error if step is not a string or a step pointer
func stepName(step any) (string, error) {
	var stepName string
	switch s := step.(type) {
	case string:
		stepName = s
	case *Step:
		stepName = s.Name
	default:
		return "", fmt.Errorf("invalid step type: %T", step)
	}
	return stepName, nil
}

// GetProgress returns the workflow progress
func (w *Workflow) GetProgress() *Progress {
	total := len(w.steps)
	completed := 0

	// Get step names
	stepNames := lo.Map(w.steps, func(item *Step, index int) string {
		return item.Name
	})

	currentStepPosition := arr.Index(stepNames, w.state.CurrentStepName)

	// Count completed steps
	for i, name := range stepNames {
		if i < currentStepPosition || w.IsStepComplete(name) {
			completed++
		}
	}

	pending := total - completed
	percents := float64(completed) / float64(total) * 100

	return &Progress{
		Total:     total,
		Completed: completed,
		Current:   currentStepPosition,
		Pending:   pending,
		Percents:  percents,
	}
}

// GetSteps returns all steps
func (w *Workflow) GetSteps() []*Step {
	return w.steps
}

// GetStep returns a step by name
func (w *Workflow) GetStep(name string) *Step {
	stepNames := lo.Map(w.steps, func(item *Step, index int) string {
		return item.Name
	})

	stepIndex := arr.Index(stepNames, name)
	if stepIndex == -1 {
		return nil
	}

	return w.steps[stepIndex]
}

// GetStepMeta returns step metadata
//
// Business logic:
// 1. Get step name
// 2. Get step metadata
// 3. Return metadata or nil if key not found
func (w *Workflow) GetStepMeta(step any, key string) any {
	stepName, err := stepName(step)
	if err != nil {
		return nil
	}

	meta, exists := w.state.StepDetails[stepName].Meta[key]
	if !exists {
		return nil
	}
	return meta
}

// SetStepMeta sets step metadata
func (w *Workflow) SetStepMeta(step any, key string, value interface{}) {
	stepName, err := stepName(step)
	if err != nil {
		return
	}

	w.state.StepDetails[stepName].Meta[key] = value
}

// MarkStepAsCompleted marks a step as completed
//
// Business logic:
// 1. Get step name
// 2. Mark step as completed
// 3. Return true if step was marked as completed
func (w *Workflow) MarkStepAsCompleted(step any) bool {
	stepName, err := stepName(step)
	if err != nil {
		return false
	}

	if _, exists := w.state.StepDetails[stepName]; !exists {
		return false
	}

	w.state.StepDetails[stepName].Completed = time.Now().Format(time.RFC3339)

	return true
}

// GetState returns the current workflow state
func (w *Workflow) GetState() *WorkflowState {
	return w.state
}

// ToString serializes the workflow state to a string
func (w *Workflow) ToString() (string, error) {
	data, err := json.Marshal(w.state)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromString deserializes the workflow state from a string
func (w *Workflow) FromString(str string) error {
	state := &WorkflowState{}
	err := json.Unmarshal([]byte(str), state)
	if err != nil {
		return err
	}
	w.state = state
	return nil
}
