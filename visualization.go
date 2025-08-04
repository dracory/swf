package swf

import (
	"bytes"
	"fmt"
	"text/template"
)

// DotNodeSpec represents a node in the DOT graph
type DotNodeSpec struct {
	Name        string
	DisplayName string
	Tooltip     string
	Shape       string
	Style       string
	FillColor   string
}

// DotEdgeSpec represents an edge in the DOT graph
type DotEdgeSpec struct {
	FromNodeName string
	ToNodeName   string
	Tooltip      string
	Style        string
	Color        string
}

const dotTemplateText = `digraph {
	rankdir = "LR"
	node [fontname="Arial"]
	edge [fontname="Arial"]
{{ range $node := $.Nodes}}	"{{$node.Name}}" [label="{{$node.DisplayName}}" shape={{$node.Shape}} style={{$node.Style}} tooltip="{{$node.Tooltip}}" fillcolor="{{$node.FillColor}}" {{if eq $node.Style "filled"}}fontcolor="white"{{end}}]
{{ end }}        
{{ range $edge := $.Edges}}	"{{$edge.FromNodeName}}" -> "{{$edge.ToNodeName}}" [style={{$edge.Style}} tooltip="{{$edge.Tooltip}}" color="{{$edge.Color}}"]
{{ end }}}`

var dotTemplate = template.Must(template.New("digraph").Parse(dotTemplateText))

// Visualize returns a DOT graph representation of the workflow
func (w *Workflow) Visualize() string {
	// Handle empty workflow
	if len(w.steps) == 0 {
		return `digraph {
	rankdir = "LR"
	node [fontname="Arial"]
	edge [fontname="Arial"]
}`
	}

	nodes := make([]*DotNodeSpec, 0, len(w.steps))
	edges := make([]*DotEdgeSpec, 0, len(w.steps)-1)

	// Create nodes
	for i, step := range w.steps {
		nodeStyle := "solid"
		fillColor := "#ffffff"

		// Current step is filled blue
		if w.IsStepCurrent(step) {
			nodeStyle = "filled"
			fillColor = "#2196F3"
		} else if w.IsStepComplete(step) {
			// Completed steps are filled green
			nodeStyle = "filled"
			fillColor = "#4CAF50"
		}

		nodes = append(nodes, &DotNodeSpec{
			Name:        step.Name,
			DisplayName: step.Title,
			Tooltip:     step.Description,
			Shape:       "box",
			Style:       nodeStyle,
			FillColor:   fillColor,
		})

		// Create edges between steps
		if i > 0 {
			edgeStyle := "solid"
			edgeColor := "#9E9E9E"

			// Highlight the path up to the current step
			if w.IsStepComplete(w.steps[i-1]) {
				edgeColor = "#4CAF50"
			}

			edges = append(edges, &DotEdgeSpec{
				FromNodeName: w.steps[i-1].Name,
				ToNodeName:   step.Name,
				Style:        edgeStyle,
				Color:        edgeColor,
				Tooltip:      fmt.Sprintf("From %s to %s", w.steps[i-1].Title, step.Title),
			})
		}
	}

	buf := new(bytes.Buffer)
	err := dotTemplate.Execute(buf, struct {
		Nodes []*DotNodeSpec
		Edges []*DotEdgeSpec
	}{
		Nodes: nodes,
		Edges: edges,
	})

	if err != nil {
		return fmt.Sprintf("Error generating DOT graph: %v", err)
	}

	return buf.String()
}
