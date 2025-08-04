package swf

// Step represents a single step in a workflow
type Step struct {
	Name        string
	Type        string
	Title       string
	Description string
	Responsible string
}

// NewStep creates a new Step with the given name
func NewStep(name string) *Step {
	return &Step{
		Name:        name,
		Type:        "normal",
		Title:       "",
		Description: "",
		Responsible: "Admin",
	}
}

// GetActionLink returns the action link for the step
// Note: This is a placeholder implementation as the PHP version uses a framework-specific function
func (s *Step) GetActionLink() string {
	// In Go, we would typically use a router or URL builder instead of the PHP action() function
	// This is a simplified implementation
	return "/" + s.Responsible + "/" + s.Name
}
