# zed
ZenHub Epic Dependencies

## Who's Zed?

zed is a tool go generate a dependency graph for zenhub epics.

```
zed uses github and zenhub apis to generate a dependency graph
        for a selected epic in a 'mermaid' md-like format.

Usage:
  zed [flags]

Flags:
  -e, --epic int          ID of an epic (required)
  -g, --github string     GitHub token (requried)
  -h, --help              help for zed
  -o, --owner string      owner of repo, must be provided if repo is referenced by name
  -r, --repo int          ID of repo (alternative '--repoName' and '--owner')
  -n, --repoName string   name of repo, requires owner to be defined (aletrnative '--repo')
  -t, --timeout int       client timeout (default 10)
  -z, --zenhub string     ZenHub token (required)
```

The output is `mermaid` markdown text which can be rendered to a flowchart of
dependencies in the epic. Each node has a link on a relevant github issue. It 
also contains some information about the issue.

----

Note there are no limitations on usage of github or zenhub API in place in zed.
If the epic you access has hundreds of issues - hundreds of request will be made.

## Build & run
```
git clone git@github.com:seriar-org/zed.git
cd zed
go build
./zed -z <ZENHUB_TOKEN> -g <GITHUB_TOKEN> -r <RepoID> -e <EpicID> > graph.md
```
You can then render results from `graph.md` in some available online tools:
* [https://mermaid.live/](https://mermaid.live/)
* [https://mermaid-js.github.io/mermaid-live-editor](https://mermaid-js.github.io/mermaid-live-editor)


## Tips

### Getting github token:
github token can be created at [https://github.com/settings/tokens](https://github.com/settings/tokens)

A token for `zed` needs `repo` scope

### Getting zenhub token
zenhub token can be created at [https://app.zenhub.com/settings/tokens](https://app.zenhub.com/settings/tokens)

Note, that only one token is supported at the time, and there are no scope limitations

### Security
Remember to be cautious not to share your tokens, as they provide read/write 
access to at least some parts of your infrastructure

### RepoID
RepositoryID can be found in zenhub board's url

`https://app.zenhub.com/workspaces/<ws-id>/board?repos=<repoIDs>`
