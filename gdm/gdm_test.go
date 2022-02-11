package gdm

import (
	"fmt"
	"testing"
)

func expectEquals(t *testing.T, message string, received string, expected string) {
	if received != expected {
		t.Error(fmt.Sprintf("%s (expected '%s', got '%s')", message, expected, received))
	}
}

func TestGraphCreation(t *testing.T) {
	expected := "graph TD;"
	m := CreateMermaidGraph()

	r := m.Render()

	expectEquals(t, "Empty graph created incorrectly", r, expected)
}

func TestAddNode(t *testing.T) {
	expected := "graph TD;\r\nnodeId(nodeText)"

	m := CreateMermaidGraph()
	m.AddNode("nodeId", "nodeText")

	r := m.Render()

	expectEquals(t, "Node added incorrectly", r, expected)
}

func TestAddLink(t *testing.T) {
	expected := "graph TD;\r\nclick nodeId \"http://link.com\" \"hint\""

	r := CreateMermaidGraph().AddLink("nodeId", "http://link.com", "hint").Render()

	expectEquals(t, "Link added incorrectly", r, expected)
}

func TestAddConnection(t *testing.T) {
	expected := "graph TD;\r\nnodeFrom --> nodeTo"

	r := CreateMermaidGraph().AddConnection("nodeFrom", "nodeTo").Render()

	expectEquals(t, "Connection added incorrectly", r, expected)
}
