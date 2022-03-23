package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/cli/go-gh"
)

type config struct {
	repo    string
	org     string
	user    string
	page    int
	verbose bool
	version bool
}

func parseFlags() config {
	repo := flag.String("repo", "", "a optional GitHub repository (i.e. 'python/peps') ; use repo for current folder if omitted and no 'org' nor 'user' flag")
	org := flag.String("org", "", "a optional GitHub organization (i.e. 'python') to scan the repositories from (100 max) ; use repo for current folder if omitted and no 'repo' nor 'user' flag")
	user := flag.String("user", "", "a optional GitHub user (i.e. 'torvalds') to scan the repositories from (100 max); use repo for current folder if omitted and no 'repo' nor 'org' flag")
	page := flag.Int("page", 1, "Page number for 'repo' and 'user' flags, 100 repositories per page")
	verbose := flag.Bool("verbose", false, "mode that outputs several lines (otherwise, outputs a one-liner) ; default: false")
	version := flag.Bool("version", false, "Output version-related information")
	flag.Parse()
	return config{*repo, *org, *user, *page, *verbose, *version}
}

type owner struct{ Login string }

type repo struct {
	Name        string
	Full_name   string
	Owner       owner
	Description string
	Topics      []string
	Visibility  string
}

type collaborator struct {
	login string
}

type version struct {
	commit string
	date   time.Time
	dirty  bool
}

func main() {
	config := parseFlags()
	if config.version {
		version := getVersion()
		dirty := ""
		if version.dirty {
			dirty = "(dirty)"
		}
		fmt.Printf("Commit %s (%s) %s\n", version.commit, version.date, dirty)
	} else if len(config.org) > 0 || len(config.user) > 0 {
		repos := []repo{}
		repos, error := getRepos(config)
		if error != nil {
			fmt.Print(error)
			os.Exit(2)
		}
		for _, repo := range repos {
			repoMessage, repo, validRepo := scanRepo(config, repo.Full_name)
			if validRepo {
				fmt.Printf(repo.Full_name + ": " + repoMessage)
				collaboratorsMessage := scanCollaborators(config, repo.Full_name)
				fmt.Printf(collaboratorsMessage)
				if strings.Compare(repo.Visibility, "public") == 0 {
					communityScoreMessage := scanCommunityScore(config, repo.Full_name)
					fmt.Printf(communityScoreMessage)
				}
			}
			println()
		}
	} else if len(config.repo) > 0 {
		repoWithOrg, error := getRepo(config)
		if error != nil {
			fmt.Print(error)
			if strings.Contains(error.Error(), "none of the git remotes configured for this repository point to a known GitHub host") {
				print("If current folder is related to a GitHub repository, please check 'gh auth status' and 'gh config list'.")
			}
			os.Exit(1)
		}
		repoMessage, repo, validRepo := scanRepo(config, repoWithOrg)
		if validRepo {
			fmt.Printf(repoMessage)
			collaboratorsMessage := scanCollaborators(config, repoWithOrg)
			fmt.Printf(collaboratorsMessage)
			if strings.Compare(repo.Visibility, "public") == 0 {
				communityScoreMessage := scanCommunityScore(config, repoWithOrg)
				fmt.Printf(communityScoreMessage)
			}
		}
	}
}

func getRepos(config config) ([]repo, error) {
	if len(config.org) == 0 && len(config.user) == 0 {
		return []repo{}, nil
	}
	client, err := gh.RESTClient(nil)
	if err != nil {
		fmt.Print(err)
		return []repo{}, err
	}
	if len(config.org) > 0 {
		// https://docs.github.com/en/rest/reference/repos#list-organization-repositories
		repos := []repo{}
		err = client.Get(
			"orgs/"+config.org+"/repos?sort=full_name&per_page=100&page="+strconv.Itoa(config.page),
			&repos)
		return repos, err
	} else {
		// https://docs.github.com/en/rest/reference/repos#list-repositories-for-a-user
		repos := []repo{}
		err = client.Get(
			"users/"+config.user+"/repos?sort=full_name&per_page=100&page="+strconv.Itoa(config.page),
			&repos)
		return repos, err
	}
}

func getRepo(config config) (string, error) {
	if len(config.repo) > 1 {
		return config.repo, nil
	}
	if config.verbose {
		fmt.Printf("(current repo)\n")
	}
	currentRepo, error := gh.CurrentRepository()
	if error != nil {
		return "", error
	}
	return currentRepo.Owner() + "/" + currentRepo.Name(), nil
}

func scanRepo(config config, repoWithOrg string) (message string, repository repo, validRepo bool) {
	// https://docs.github.com/en/rest/reference/repos#get-a-repository-readme
	readme := struct {
		Name string
	}{}
	client, err := gh.RESTClient(nil)
	if err != nil {
		fmt.Print(err)
		return
	}
	err = client.Get(
		"repos/"+repoWithOrg+"/readme",
		&readme)
	if len(readme.Name) > 0 {
		if config.verbose {
			message = "\n" + message + "  - a README ☑️\n"
		} else {
			message = message + "README ☑️, "
		}
	} else if strings.HasPrefix(err.Error(), "HTTP 404: Not Found") {
		if config.verbose {
			message = "\n" + message + "  - no README 😇, \n"
		} else {
			message = message + "no README 😇, "
		}
	} else {
		fmt.Print(err)
	}

	repo := struct {
		Name        string
		Full_name   string
		Owner       owner
		Description string
		Topics      []string
		Visibility  string
	}{}
	errRepo := client.Get(
		"repos/"+repoWithOrg,
		&repo)
	if errRepo != nil {
		fmt.Print(errRepo)
		return
	}
	if len(repo.Topics) > 0 {
		if config.verbose {
			message = message + "  - topics ☑️\n"
		} else {
			message = message + "topics ☑️, "
		}
	} else {
		if config.verbose {
			message = message + "  - no topics 😇\n"
		} else {
			message = message + "no topics 😇, "
		}
	}
	return message, repo, true
}

func scanCollaborators(config config, repoWithOrg string) string {
	// https://docs.github.com/en/rest/reference/collaborators#list-repository-collaborators
	client, err := gh.RESTClient(nil)
	if err != nil {
		fmt.Print(err)
		return ""
	}
	collaborators := []collaborator{}
	err = client.Get(
		"repos/"+repoWithOrg+"/collaborators",
		&collaborators)
	message := ""
	if err != nil && len(err.Error()) > 0 {
		if strings.HasPrefix(err.Error(), "HTTP 403") {
			// 🤫
		} else {
			fmt.Print(err)
		}
	} else if len(collaborators) <= 1 {
		if config.verbose {
			message = message + fmt.Sprintf("  - %d collaborator 👤", len(collaborators))
		} else {
			message = message + fmt.Sprintf("%d collaborator 👤 ", len(collaborators))
		}
	} else {
		if config.verbose {
			message = message + fmt.Sprintf("  - %d collaborators 👥", len(collaborators))
		} else {
			message = message + fmt.Sprintf("%d collaborators 👥 ", len(collaborators))
		}
	}
	return message
}

func scanCommunityScore(config config, repoWithOrg string) string {
	// https://docs.github.com/en/rest/reference/metrics#get-community-profile-metrics
	communityProfile := struct {
		Health_percentage int64
	}{}
	client, err := gh.RESTClient(nil)
	if err != nil {
		fmt.Print(err)
		return ""
	}
	err = client.Get(
		"repos/"+repoWithOrg+"/community/profile",
		&communityProfile)
	if err != nil {
		fmt.Print(err)
		return ""
	}
	message := ""
	if config.verbose {
		message = message + fmt.Sprintf("  - a community profile score of %d 💯", communityProfile.Health_percentage)
	} else {
		message = message + fmt.Sprintf("community profile score: %d 💯", communityProfile.Health_percentage)
	}
	return message
}

func getVersion() version {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		panic("Cannot read build info")
	}
	revision := "?"
	dirtyBuild := false
	date := time.Now()
	for _, kv := range info.Settings {
		switch kv.Key {
		case "vcs.revision":
			revision = kv.Value
		case "vcs.time":
			date, _ = time.Parse(time.RFC3339, kv.Value)
		case "vcs.modified":
			dirtyBuild = kv.Value == "true"
		}
	}
	return version{revision, date, dirtyBuild}
}
