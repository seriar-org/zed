package gdm

import (
	"fmt"
	"strings"
)

type Mermaid struct {
	data []string
}

func CreateMermaidGraph() *Mermaid {
	data := make([]string, 0)
	m := &Mermaid{
		data: data,
	}
	m.Graph("TD")
	return m
}

func (m *Mermaid) Render() string {
	return strings.Join(m.data, "\r\n")
}

func (m *Mermaid) AddLine(line string) *Mermaid {
	m.data = append(m.data, line)
	return m
}

func (m *Mermaid) Graph(orientation string) *Mermaid {
	m.AddLine(fmt.Sprintf("graph %s;", orientation))
	return m
}

func (m *Mermaid) AddNode(id string, text string) *Mermaid {
	m.AddLine(fmt.Sprintf("%s(%s)", id, text))
	return m
}

func (m *Mermaid) AddLink(id string, link string, hint string) *Mermaid {
	m.AddLine(fmt.Sprintf("click %s \"%s\" \"%s\"", id, link, hint))
	return m
}

func (m *Mermaid) AddConnection(from string, to string) *Mermaid {
	m.AddLine(fmt.Sprintf("%s --> %s", from, to))
	return m
}
