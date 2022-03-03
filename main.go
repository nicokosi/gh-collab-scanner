package main

import (
	"flag"
	"fmt"
	"os"
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
	repoWithOrg, error := getRepo(config)
	if error != nil {
		fmt.Print(error)
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
		fmt.Println(err)
		return
	}
	err = client.Get(
		"repos/"+repoWithOrg+"/readme",
		&readme)
	if len(readme.Name) > 0 {
		if config.verbose {
			message = message + "  - a README ☑️\n"
		} else {
			message = message + "README ☑️, "
		}
	} else if strings.HasPrefix(err.Error(), "HTTP 404: Not Found") {
		if config.verbose {
			message = message + "no README 😇, \n"
		} else {
			message = message + "no README 😇, "
		}
	} else {
		fmt.Println(err)
	}

	repo := struct {
		Name        string
		Owner       owner
		Description string
		Topics      []string
		Visibility  string
	}{}
	errRepo := client.Get(
		"repos/"+repoWithOrg,
		&repo)
	if errRepo != nil {
		fmt.Println(errRepo)
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
		fmt.Println(err)
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
			fmt.Println(err)
		}
	} else if len(collaborators) <= 1 {
		if config.verbose {
			message = message + fmt.Sprintf("  - %d collaborator 👤\n", len(collaborators))
		} else {
			message = message + fmt.Sprintf("%d collaborator 👤, ", len(collaborators))
		}
	} else {
		if config.verbose {
			message = message + fmt.Sprintf("  - %d collaborators 👥\n", len(collaborators))
		} else {
			message = message + fmt.Sprintf("%d collaborators 👥, ", len(collaborators))
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
		fmt.Println(err)
		return ""
	}
	err = client.Get(
		"repos/"+repoWithOrg+"/community/profile",
		&communityProfile)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	message := ""
	if config.verbose {
		message = message + fmt.Sprintf("  - a community profile score of %d 💯\n", communityProfile.Health_percentage)
	} else {
		message = message + fmt.Sprintf("community profile score: %d 💯\n", communityProfile.Health_percentage)
	}
	return message
}
