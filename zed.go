package main

import (
	"fmt"

	"github.com/seriar-org/zed/gdm"
	"github.com/seriar-org/zed/gzc"
)

func CreateIssueNodes(repoID, epicID int, client *gzc.Client, graph *gdm.Mermaid) (*gdm.Mermaid, error) {

	e, err := client.RequestEpic(repoID, epicID)
	if err != nil {
		return nil, err
	}

	for index, issue := range e.Issues {
		issueID := fmt.Sprintf("issue_%d", index)
		text := fmt.Sprintf("repo-%d#issue-%d", issue.RepoID, issue.IssueNumber)
		graph.AddNode(issueID, text)
	}

	return graph, nil

}
