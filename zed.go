package main

import (
	"fmt"

	"github.com/seriar-org/zed/gdm"
	"github.com/seriar-org/zed/gzc"
)

func createIssueID(i gzc.Issue) string {
	return fmt.Sprintf("r%di%d", i.RepoID, i.IssueNumber)
}
func CreateIssueNodes(repoID, epicID int, client *gzc.Client, graph *gdm.Mermaid) (*gdm.Mermaid, *gzc.Epic, error) {

	e, err := client.RequestEpic(repoID, epicID)
	if err != nil {
		return nil, nil, err
	}

	for _, issue := range e.Issues {
		issueID := createIssueID(issue)
		text := fmt.Sprintf("repo-%d#issue-%d", issue.RepoID, issue.IssueNumber)
		graph.AddNode(issueID, text)
	}

	return graph, e, nil

}

func CreateDependencyLinks(epic *gzc.Epic, client *gzc.Client, graph *gdm.Mermaid) (*gdm.Mermaid, error) {
	// creating epic repo-issues multimap
	issuesByRepo := make(map[int][]int)
	for _, issue := range epic.Issues {
		issuesByRepo[issue.RepoID] = append(issuesByRepo[issue.RepoID], issue.IssueNumber)
	}

	// for each repo in epic getting all dependencies
	for repoID, _ := range issuesByRepo {
		dependencies, err := client.RequestDependencies(repoID)
		if err != nil {
			return nil, err
		}
		// for each dependency checking if blocked issue is in epic
		// TODO can probably use multimap
		for _, dependency := range dependencies.Dependencies {
			if epic.HasIssue(dependency.Blocked.RepoID, dependency.Blocked.IssueNumber) {
				target := createIssueID(dependency.Blocked)
				source := createIssueID(dependency.Blocking)

				// adding a connection
				graph.AddConnection(source, target)
			}
		}
	}
	return graph, nil
}
