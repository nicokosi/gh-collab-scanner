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
