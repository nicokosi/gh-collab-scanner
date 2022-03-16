# `collab-scanner` GitHub CLI extension

A [GitHub CLI extension](https://docs.github.com/en/github-cli/github-cli/using-github-cli-extensions) that displays collaboration-related information on repositories.

![collab-scanner](https://user-images.githubusercontent.com/3862051/157172870-0d50c1b8-d238-4227-ad86-d12855303e13.gif)

## Install ⬇️

```sh
gh extension install nicokosi/gh-collab-scanner
```

## Use ▶️

From a folder where a GitHub repository has been cloned:

```sh
gh collab-scanner
```

will display something like:

> (current repo)
> Repo nicokosi/gh-collab-scanner has: description ☑️, README ☑️, topics ☑️, 1 collaborator 👤, community profile score: 33 💯

For any GitHub repository, via its full name ${org}/${repo} (i.e. python/peps), use the `--repo` flag:

```sh
gh collab-scanner --repo python/peps
```

It will display something like:

> Repo python/peps has: description ☑️, README ☑️, no topics 😇, community profile score: 71 💯

In order to scan all repositories for a given organization, use the `--org` flag:

```sh
gh collab-scanner --org python
```

In order to scan all repositories for a given user, use the `--user` flag:

```sh
gh collab-scanner --user torvalds
```

Need help? Run:

```sh
gh-collab-scanner --help
```

## Develop 🧑‍💻

### Build from source code 🧑‍💻▶️

Build then run:

```sh
go build && ./gh-collab-scanner
```

### Install from source code 🧑‍💻⏺

Install and run:

```sh
gh extension install .; gh collab-scanner
```

### Examine code 🧑‍💻🔍

```sh
go vet
```

### Run tests 🧑‍💻☑️

```sh
go test -v -cover
```

### Release 🧑‍💻📦

Check the current version:

```sh
gh release view | head -n 2
```

Then create a tag for the next version with respect with [semver](https://semver.org):

```sh
git tag ${version}
git push origin ${version}
```
