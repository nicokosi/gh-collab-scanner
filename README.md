# `collab-scanner` GitHub CLI extension

A [GitHub CLI extension](https://docs.github.com/en/github-cli/github-cli/using-github-cli-extensions) that displays collaboration-related information on a repository.

## Install

```sh
gh extension install nicokosi/gh-collab-scanner
```

## Use

From a folder where a GitHub repository has been cloned:

```sh
gh collab-scanner
```

will display something like:

  (current repo)
  Repo nicokosi/gh-collab-scanner has: description ☑️, README ☑️, topics ☑️, 1 collaborator 👤, community profile score: 33 💯

For any GitHub repository via its full name `org`/`repo` (i.e. python/peps)

```sh
gh collab-scanner --repo python/peps
```

will display something like:

  Repo python/peps has: description ☑️, README ☑️, no topics 😇, community profile score: 71 💯

Need help? Run:

```sh
gh-collab-scanner --help
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
