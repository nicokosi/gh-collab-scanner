# `collab-scanner` GitHub CLI extension

A [GitHub CLI extension](https://docs.github.com/en/github-cli/github-cli/using-github-cli-extensions) that displays collaboration-related information on a repository:

```sh
# For current repository
$ gh collab-scanner
(current repo)
Repository nicokosi/gh-collab-scanner has:
  - a description ☑️
  - no README 😇
  - no topics 😇
  - 1 collaborator 👤
  - a community profile score of 16 💯
```

```sh
# For a given repository, using its full name ${organization}/${name}
$ gh collab-scanner facebook/react
Repository facebook/react has:
  - a description ☑️
  - has a README ☑️
  - topics ☑️
  - a community profile score of 85 💯
```

## Build/install from source code

Build an run:

```sh
go build && ./gh-collab-scanner
```

Install and run:

```sh
gh extension install .; gh collab-scanner
```
