package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v42/github"
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
	Repo        string
	Owner       string
	IssueID     int
	StoryPoints int
	State       string
	Title       string
	Assignee    string
	URL         string
}

func (n *Node) MermaidNodeText() string {
	return fmt.Sprintf("%s/%s#%d<br/>%s<br/>%d SP<br/>%s<br/>%s", n.Owner, n.Repo, n.IssueID, n.Title, n.StoryPoints, n.Assignee, n.State)
}

type Edge struct {
	SourceID string
	TargetID string
}

type RepoWithOwner struct {
	Owner      string
	Repository string
}

type Zed struct {
	DependenciesGraph DependenciesGraph
	Repos             map[int]RepoWithOwner
	Mermaid           *gdm.Mermaid
	ZenHubClient      *gzc.Client
	GitHubClient      *github.Client
	Context           context.Context
	Epic              *gzc.Epic
}

func CreateZed(zenhubClient *gzc.Client, githubClient *github.Client, context context.Context) *Zed {
	return &Zed{
		DependenciesGraph: DependenciesGraph{
			Nodes: make(map[string]Node, 0),
			Edges: make([]Edge, 0),
		},
		Repos:        make(map[int]RepoWithOwner),
		Mermaid:      gdm.CreateMermaidGraph(),
		ZenHubClient: zenhubClient,
		GitHubClient: githubClient,
		Context:      context,
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

		repoWithOwner, ok := z.Repos[issue.RepoID]
		if !ok {
			r, _, err := z.GitHubClient.Repositories.GetByID(z.Context, int64(issue.RepoID))
			if err != nil {
				return nil, err
			}
			repoWithOwner = RepoWithOwner{*r.Owner.Login, *r.Name}
			z.Repos[issue.RepoID] = repoWithOwner
		}

		owner := repoWithOwner.Owner
		repo := repoWithOwner.Repository

		i, _, err := z.GitHubClient.Issues.Get(z.Context, owner, repo, issue.IssueNumber)
		if err != nil {
			return nil, err
		}
		assignee := "unassigned"

		if i.Assignee != nil {
			assignee = *i.Assignee.Login
		}

		n := Node{
			Repo:        repo,
			Owner:       owner,
			IssueID:     issue.IssueNumber,
			State:       *i.State,
			Assignee:    assignee,
			Title:       *i.Title,
			URL:         *i.HTMLURL,
			StoryPoints: issue.Estimate.Value,
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
	for repoID := range issuesByRepo {
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
		z.Mermaid.AddLink(id, node.URL, "on github")
	}
	for _, edge := range z.DependenciesGraph.Edges {
		z.Mermaid.AddConnection(edge.SourceID, edge.TargetID)
	}
	return z.Mermaid.Render()
}
