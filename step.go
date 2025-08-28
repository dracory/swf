package swf

// Step represents a single step in a workflow
type Step struct {
	// Name is a unique identifier for the step.
	// Used internally for workflow logic and step identification.
	// Must be unique within a workflow and should be in snake_case (e.g., 'document_review').
	Name string

	// Type defines the step's behavior or category.
	// Common types include 'normal', 'approval', 'notification', etc.
	// Defaults to 'normal' if not specified.
	Type string

	// Title is the human-readable display name of the step.
	// Used in UI/display purposes and can contain spaces and special characters.
	// Example: 'Document Review' or 'Manager Approval'.
	Title string

	// Description provides detailed information about what the step does.
	// This is used for documentation and can be displayed in the UI to help users understand the step's purpose.
	// Example: 'Review and approve the document before it proceeds to the next stage.'
	Description string

	// Responsible identifies who is accountable for completing this step.
	// Can be a role (e.g., 'manager'), team name, or specific person's identifier.
	// Used for routing and notification purposes.
	// Example: 'admin@example.com' or 'document_approvers'
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
