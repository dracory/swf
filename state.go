package swf

// State represents a workflow state
type State struct {
	Name        string
	Title       string
	Description string
	Responsible string
}

// NewState creates a new State with the given name
func NewState(name string) *State {
	return &State{
		Name:        name,
		Title:       "",
		Description: "",
		Responsible: "Admin",
	}
}

// GetActionLink returns the action link for the state
func (s *State) GetActionLink() string {
	return "/" + s.Responsible + "/" + s.Name
}
