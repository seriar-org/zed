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
	Edges map[Edge]struct{} // set
}

func (g *DependenciesGraph) HasIssue(issue gzc.Issue) bool {
	_, ok := g.Nodes[createIssueID(issue)]
	return ok
}

type Node struct {
	Repo        string
	Owner       string
	IssueID     int
	StoryPoints float64
	State       string
	Title       string
	Assignee    string
	URL         string
	IsExternal  bool
}

func (n *Node) MermaidNodeText() string {
	sp := "?"
	if n.StoryPoints > 0 {
		sp = fmt.Sprintf("%.1f", n.StoryPoints)
	}
	text := fmt.Sprintf("%s/%s#%d<br/>%s<br/>%s SP<br/>%s<br/>%s", n.Owner, n.Repo, n.IssueID, n.Title, sp, n.Assignee, n.State)
	if n.IsExternal {
		text = fmt.Sprintf("<i>%s<br/>EXTERNAL</i>", text)
	}
	return text
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

func CreateZed(context context.Context, zenhubClient *gzc.Client, githubClient *github.Client) *Zed {
	return &Zed{
		DependenciesGraph: DependenciesGraph{
			Nodes: make(map[string]Node, 0),
			Edges: make(map[Edge]struct{}, 0),
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

func (z *Zed) createNode(repoID, issueID int) (*Node, error) {
	repoWithOwner, ok := z.Repos[repoID]
	if !ok {
		r, _, err := z.GitHubClient.Repositories.GetByID(z.Context, int64(repoID))
		if err != nil {
			return nil, err
		}
		repoWithOwner = RepoWithOwner{*r.Owner.Login, *r.Name}
		z.Repos[repoID] = repoWithOwner
	}

	owner := repoWithOwner.Owner
	repo := repoWithOwner.Repository

	i, _, err := z.GitHubClient.Issues.Get(z.Context, owner, repo, issueID)
	if err != nil {
		return nil, err
	}
	assignee := "unassigned"

	if i.Assignee != nil {
		assignee = *i.Assignee.Login
	}

	return &Node{
		Repo:     repo,
		Owner:    owner,
		IssueID:  issueID,
		State:    *i.State,
		Assignee: assignee,
		Title:    *i.Title,
		URL:      *i.HTMLURL,
	}, nil
}

func (z *Zed) CreateIssueNodes(repoID, epicID int) (*Zed, error) {

	e, err := z.ZenHubClient.RequestEpic(repoID, epicID)
	if err != nil {
		return nil, err
	}
	z.Epic = e

	for _, issue := range e.Issues {
		n, err := z.createNode(issue.RepoID, issue.IssueNumber)
		if err != nil {
			return nil, err
		}
		n.StoryPoints = issue.Estimate.Value
		nodeID := createIssueID(issue)
		z.DependenciesGraph.Nodes[nodeID] = *n
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
		for _, dependency := range dependencies.Dependencies {
			targetID := createIssueID(dependency.Blocked)
			if z.DependenciesGraph.HasIssue(dependency.Blocked) {
				sourceID := createIssueID(dependency.Blocking)
				if _, ok := z.DependenciesGraph.Nodes[sourceID]; !ok {

					n, err := z.createNode(dependency.Blocking.RepoID,
						dependency.Blocking.IssueNumber)
					if err != nil {
						return nil, err
					}
					n.StoryPoints = -1
					n.IsExternal = true
					z.DependenciesGraph.Nodes[sourceID] = *n
				}
				// adding a connection
				z.DependenciesGraph.Edges[Edge{sourceID, targetID}] = struct{}{} //add to set
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
	for edge := range z.DependenciesGraph.Edges {
		z.Mermaid.AddConnection(edge.SourceID, edge.TargetID)
	}
	return z.Mermaid.Render()
}
