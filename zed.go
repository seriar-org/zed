package main

import (
	"fmt"

	"github.com/seriar-org/zed/gdm"
	"github.com/seriar-org/zed/gzc"
)

type DependenciesGraph struct {
	Nodes map[string]Node
	Edges []Edge
}

func (g *DependenciesGraph) HasIssue(issue gzc.Issue) bool {
	_, ok := g.Nodes[createIssueID(issue)]
	return ok
}

type Node struct {
	GitHubLink  string
	Repo        string
	IssueID     int
	StoryPoints int
	Status      string
}

func (n *Node) MermaidNodeText() string {
	return fmt.Sprintf("%s#%d", n.Repo, n.IssueID)
}

type Edge struct {
	SourceID string
	TargetID string
}

type Zed struct {
	DependenciesGraph DependenciesGraph
	Mermaid           *gdm.Mermaid
	ZenHubClient      *gzc.Client
	Epic              *gzc.Epic
}

func CreateZed(client *gzc.Client) *Zed {
	return &Zed{
		DependenciesGraph: DependenciesGraph{
			Nodes: make(map[string]Node, 0),
			Edges: make([]Edge, 0),
		},
		Mermaid:      gdm.CreateMermaidGraph(),
		ZenHubClient: client,
	}
}

func createIssueID(i gzc.Issue) string {
	return fmt.Sprintf("r%di%d", i.RepoID, i.IssueNumber)
}

func (z *Zed) CreateIssueNodes(repoID, epicID int) (*Zed, error) {

	e, err := z.ZenHubClient.RequestEpic(repoID, epicID)
	if err != nil {
		return nil, err
	}
	z.Epic = e

	for _, issue := range e.Issues {
		n := Node{
			Repo:    fmt.Sprintf("repo/%d", issue.RepoID),
			IssueID: issue.IssueNumber,
		}
		issueID := createIssueID(issue)
		z.DependenciesGraph.Nodes[issueID] = n
	}

	return z, nil

}

func (z *Zed) CreateDependencyLinks() (*Zed, error) {
	// creating epic repo-issues multimap
	issuesByRepo := make(map[int][]int)
	for _, issue := range z.Epic.Issues {
		issuesByRepo[issue.RepoID] = append(issuesByRepo[issue.RepoID], issue.IssueNumber)
	}

	// for each repo in epic getting all dependencies
	for repoID, _ := range issuesByRepo {
		dependencies, err := z.ZenHubClient.RequestDependencies(repoID)
		if err != nil {
			return nil, err
		}
		// for each dependency checking if blocked issue is in epic
		// TODO can probably use multimap
		for _, dependency := range dependencies.Dependencies {
			if z.DependenciesGraph.HasIssue(dependency.Blocked) {
				target := createIssueID(dependency.Blocked)
				source := createIssueID(dependency.Blocking)

				// adding a connection
				z.DependenciesGraph.Edges = append(z.DependenciesGraph.Edges, Edge{source, target})
			}
		}
	}
	return z, nil
}

func (z *Zed) Render() string {
	for id, node := range z.DependenciesGraph.Nodes {
		z.Mermaid.AddNode(id, node.MermaidNodeText())
	}
	for _, edge := range z.DependenciesGraph.Edges {
		z.Mermaid.AddConnection(edge.SourceID, edge.TargetID)
	}
	return z.Mermaid.Render()
}
