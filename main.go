package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/google/go-github/v42/github"
	"github.com/seriar-org/zed/gzc"
	"golang.org/x/oauth2"
)

func parseArgs() (string, string, int, int, int) {
	zenhubToken := os.Args[1]
	githubToken := os.Args[2]
	repoID, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic("Cannot convert repository id to int")
	}
	epicID, err := strconv.Atoi(os.Args[4])
	if err != nil {
		panic("Cannot convert epic id to int")
	}
	timeout, err := strconv.Atoi(os.Args[5])
	if err != nil {
		panic("Cannot convert timeout to int")
	}
	return zenhubToken, githubToken, repoID, epicID, timeout
}

func createClient(token string, timeout int) *gzc.Client {
	a := gzc.CreateAPI(&http.Client{}, "https://api.zenhub.com").WithTimeout(timeout)
	c := gzc.CreateClient(a, token)
	return c
}

func main() {
	fmt.Println("Who's Zed?")

	zenhubToken, githubToken, repoID, epicID, timeout := parseArgs()
	zenhub := createClient(zenhubToken, timeout)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	github := github.NewClient(tc)

	z := CreateZed(ctx, zenhub, github)

	_, err := z.CreateIssueNodes(repoID, epicID)
	if err != nil {
		panic(err)
	}

	_, err = z.CreateDependencyLinks()
	if err != nil {
		panic(err)
	}

	fmt.Printf("```mermaid\n%s\n```\n", z.Render())
}
