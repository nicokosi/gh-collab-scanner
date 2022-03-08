# `collab-scanner` GitHub CLI extension

A [GitHub CLI extension](https://docs.github.com/en/github-cli/github-cli/using-github-cli-extensions) that displays collaboration-related information on a repository.

![Kapture 2022-02-26 at 16 26 46](https://user-images.githubusercontent.com/3862051/155848808-716c27ba-7566-447c-8e2b-e1b2601f37a5.gif)


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

For any GitHub repository, via its full name ${org}/${repo} (i.e. python/peps):

```sh
gh collab-scanner --repo python/peps
```

will display something like:

  Repo python/peps has: description ☑️, README ☑️, no topics 😇, community profile score: 71 💯

Need help? Run:

```sh
gh-collab-scanner --help
```

## Develop

### Build from source code ▶️

Build then run:

```sh
go build && ./gh-collab-scanner
```

### Install from source code ⏺

Install and run:

```sh
gh extension install .; gh collab-scanner
```

### Examine code 🔍

```sh
go vet
```

### Run tests ☑️

```sh
go test -v -cover
```

### Release 📦

Check the current version:

```sh
gh release view | head -n 2
```

Then create a tag for the next version with respect with [semver](https://semver.org):

```sh
git tag ${version}
git push origin ${version}
```
