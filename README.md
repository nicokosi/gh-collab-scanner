# `collab-scanner` GitHub CLI extension

A [GitHub CLI extension](https://docs.github.com/en/github-cli/github-cli/using-github-cli-extensions) that displays collaboration-related information on a repository.

![gh-collab-scanner-small](https://user-images.githubusercontent.com/3862051/155272593-7ff4d205-3e0d-44df-a035-57b36a50b98a.gif)


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
