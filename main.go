package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v42/github"
	"github.com/seriar-org/zed/gzc"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

type Conf struct {
	githubToken string
	zenhubToken string
	repoName    string
	owner       string
	repoID      int
	epicID      int
	timeout     int
}

func parseArgs() (*Conf, error) {
	var (
		githubToken, zenhubToken, repoName, owner string
		repoID, epicID, timeout                   int
	)

	rootCmd := &cobra.Command{
		Use:   "zed",
		Short: "zed is zenhub epic dependencies visualizer",
		Long: `zed uses github and zenhub apis to generate a dependency graph 
	for a selected epic in a 'mermaid' md-like format.`,
		Run: func(c *cobra.Command, args []string) {},
	}
	rootCmd.Flags().StringVarP(&githubToken, "github", "g", "", "GitHub token (requried)")
	rootCmd.Flags().StringVarP(&zenhubToken, "zenhub", "z", "", "ZenHub token (required)")
	rootCmd.Flags().IntVarP(&repoID, "repo", "r", 0, "ID of repo (alternative '--repoName' and '--owner')")
	rootCmd.Flags().StringVarP(&repoName, "repoName", "n", "", "name of repo, requires owner to be defined (aletrnative '--repo')")
	rootCmd.Flags().StringVarP(&owner, "owner", "o", "", "owner of repo, must be provided if repo is referenced by name")
	rootCmd.Flags().IntVarP(&epicID, "epic", "e", 0, "ID of an epic (required)")
	rootCmd.Flags().IntVarP(&timeout, "timeout", "t", 10, "client timeout")
	err := rootCmd.Execute()

	if err != nil {
		rootCmd.Help()
		return nil, err
	}
	if githubToken == "" || zenhubToken == "" {
		rootCmd.Help()
		return nil, fmt.Errorf("Both zenhub and github tokens must be provided. Tokens available: Github - %v, Zenhub - %v", githubToken != "", zenhubToken != "")
	}
	if repoID > 0 && repoName != "" {
		rootCmd.Help()
		return nil, fmt.Errorf("Cannot define repo by id(%d) and name(%s) both (using id is preferred)", repoID, repoName)
	}
	if owner == "" && repoName != "" {
		rootCmd.Help()
		return nil, fmt.Errorf("When defining repo by name(%s) owner must be specified", repoName)
	}
	if repoID <= 0 && repoName == "" {
		rootCmd.Help()
		return nil, fmt.Errorf("Repo must be defined either by positive id(was %d) or by non empty owner and name", repoID)
	}
	if epicID <= 0 {
		rootCmd.Help()
		return nil, fmt.Errorf("Epic must be defined by positive id (was %d)", epicID)
	}
	if timeout < 1 {
		rootCmd.Help()
		return nil, fmt.Errorf("Timeout(%d) must be more than 0", timeout)
	}
	return &Conf{
		zenhubToken: zenhubToken,
		githubToken: githubToken,
		repoName:    repoName,
		owner:       owner,
		repoID:      repoID,
		epicID:      epicID,
		timeout:     timeout}, nil
}

func createClient(token string, timeout int) *gzc.Client {
	a := gzc.CreateAPI(&http.Client{}, "https://api.zenhub.com").WithTimeout(timeout)
	c := gzc.CreateClient(a, token)
	return c
}

func main() {

	conf, err := parseArgs()
	if err != nil {
		panic(err)
	}

	zenhub := createClient(conf.zenhubToken, conf.timeout)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: conf.githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	github := github.NewClient(tc)

	z := CreateZed(ctx, zenhub, github)

	if conf.repoID <= 0 {
		r, _, err := github.Repositories.Get(ctx, conf.owner, conf.repoName)
		if err != nil {
			panic(err)
		}
		conf.repoID = int(*r.ID)
	}

	_, err = z.CreateIssueNodes(conf.repoID, conf.epicID)
	if err != nil {
		panic(err)
	}

	_, err = z.CreateDependencyLinks()
	if err != nil {
		panic(err)
	}

	fmt.Println(z.Render())
}
