package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/cli/go-gh"
)

type config struct {
	repo    string
	verbose bool
}

func parseFlags() config {
	repo := flag.String("repo", "", "a optional GitHub repository (i.e. 'python/peps') ; use repo for current folder if omitted")
	verbose := flag.Bool("verbose", false, "mode that outputs several lines (otherwise, outputs a one-liner) ; default: false")
	flag.Parse()
	return config{*repo, *verbose}
}

type owner struct{ Login string }

type repo struct {
	Name        string
	Owner       owner
	Description string
	Topics      []string
	Visibility  string
}

type collaborator struct {
	login string
}

func main() {
	config := parseFlags()
	repoWithOrg := getRepo(config)
	repo := printRepo(config, repoWithOrg)
	printCollaborators(config, repoWithOrg)
	if strings.Compare(repo.Visibility, "public") == 0 {
		printCommunityScore(config, repoWithOrg)
	}
}

func getRepo(config config) string {
	if len(config.repo) > 1 {
		return config.repo
	}
	if config.verbose {
		fmt.Printf("(current repo)\n")
	}
	currentRepo, _ := gh.CurrentRepository()
	return currentRepo.Owner() + "/" + currentRepo.Name()
}

func printRepo(config config, repoWithOrg string) repo {
	// https://docs.github.com/en/rest/reference/repos#get-a-repository-readme
	readme := struct {
		Name string
	}{}
	client, errReadme := gh.RESTClient(nil)
	errReadme = client.Get(
		"repos/"+repoWithOrg+"/readme",
		&readme)

	if len(readme.Name) > 0 {
		if config.verbose {
			fmt.Printf("  - a README ☑️\n")
		} else {
			fmt.Printf("README ☑️, ")
		}
	} else if strings.HasPrefix(errReadme.Error(), "HTTP 404: Not Found") {
		if config.verbose {
			fmt.Printf("no README 😇, \n")
		} else {
			fmt.Printf("no README 😇, ")
		}
	} else {
		fmt.Println(errReadme)
	}

	repo := struct {
		Name        string
		Owner       owner
		Description string
		Topics      []string
		Visibility  string
	}{}
	client, errRepo := gh.RESTClient(nil)
	errRepo = client.Get(
		"repos/"+repoWithOrg,
		&repo)
	if errRepo != nil {
		fmt.Println(errRepo)
	}
	if len(repo.Topics) > 0 {
		if config.verbose {
			fmt.Printf("  - topics ☑️\n")
		} else {
			fmt.Printf("topics ☑️, ")
		}
	} else {
		if config.verbose {
			fmt.Printf("  - no topics 😇\n")
		} else {
			fmt.Printf("no topics 😇, ")
		}
	}
	return repo
}

func printCollaborators(config config, repoWithOrg string) []collaborator {
	// https://docs.github.com/en/rest/reference/collaborators#list-repository-collaborators
	client, errCollabs := gh.RESTClient(nil)
	collaborators := []collaborator{}
	errCollabs = client.Get(
		"repos/"+repoWithOrg+"/collaborators",
		&collaborators)
	if errCollabs != nil && len(errCollabs.Error()) > 0 {
		if strings.HasPrefix(errCollabs.Error(), "HTTP 403") {
			// 🤫
		} else {
			fmt.Println(errCollabs)
		}
	} else if len(collaborators) <= 1 {
		if config.verbose {
			fmt.Printf("  - %d collaborator 👤\n", len(collaborators))
		} else {
			fmt.Printf("%d collaborator 👤, ", len(collaborators))
		}
	} else {
		if config.verbose {
			fmt.Printf("  - %d collaborators 👥\n", len(collaborators))
		} else {
			fmt.Printf("%d collaborators 👥, ", len(collaborators))
		}
	}
	return collaborators
}

func printCommunityScore(config config, repoWithOrg string) {
	// https://docs.github.com/en/rest/reference/metrics#get-community-profile-metrics
	communityProfile := struct {
		Health_percentage int64
	}{}
	client, errCommunityProfile := gh.RESTClient(nil)
	errCommunityProfile = client.Get(
		"repos/"+repoWithOrg+"/community/profile",
		&communityProfile)
	if errCommunityProfile != nil {
		fmt.Println(errCommunityProfile)
	}
	if config.verbose {
		fmt.Printf("  - a community profile score of %d 💯\n", communityProfile.Health_percentage)
	} else {
		fmt.Printf("community profile score: %d 💯\n", communityProfile.Health_percentage)
	}
}
