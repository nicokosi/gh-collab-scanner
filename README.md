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

```
(current repo)
Repository nicokosi/gh-collab-scanner has:
  - a description ☑️
  - no README 😇
  - no topics 😇
  - 1 collaborator 👤
  - a community profile score of 16 💯
```

For any GitHub repository via its full name `org`/`repo` (i.e. python/peps)

```sh
gh collab-scanner python/peps
```

will display something like:

```
Repository python/peps has:
  - a description ☑️
  - has a README ☑️
  - no topics 😇
  - a community profile score of 71 💯
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
