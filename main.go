package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/seriar-org/zed/gdm"
	"github.com/seriar-org/zed/gzc"
)

func parseArgs() (string, int, int, int) {
	token := os.Args[1]
	repoID, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic("Cannot convert repository id to int")
	}
	epicID, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic("Cannot convert epic id to int")
	}
	timeout, err := strconv.Atoi(os.Args[4])
	if err != nil {
		panic("Cannot convert timeout to int")
	}
	return token, repoID, epicID, timeout
}

func createClient(token string, timeout int) *gzc.Client {
	a := gzc.CreateAPI(&http.Client{}, "https://api.zenhub.com").WithTimeout(timeout)
	c := gzc.CreateClient(a, token)
	return c
}

func main() {
	fmt.Println("Who's Zed?")

	token, repoID, epicID, timeout := parseArgs()
	c := createClient(token, timeout)
	graph := gdm.CreateMermaidGraph()

	graph, e, err := CreateIssueNodes(repoID, epicID, c, graph)
	if err != nil {
		panic(err)
	}
	graph, err = CreateDependencyLinks(e, c, graph)
	if err != nil {
		panic(err)
	}

	fmt.Printf("```mermaid\n%s\n```\n", graph.Render())
}
