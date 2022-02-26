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

func main() {
	var repoWithOrg = ""
	config := parseFlags()
	if len(config.repo) > 1 {
		repoWithOrg = config.repo
	} else {
		if config.verbose {
			fmt.Printf("(current repo)\n")
		}
		currentRepo, _ := gh.CurrentRepository()
		repoWithOrg = currentRepo.Owner() + "/" + currentRepo.Name()
	}

	client, err := gh.RESTClient(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	user := struct{ Login string }{}
	err = client.Get("user", &user)

	// https://docs.github.com/en/rest/reference/repos#get-a-repository
	type Owner struct{ Login string }
	repo := struct {
		Name        string
		Owner       Owner
		Description string
		Topics      []string
		Visibility  string
	}{}
	client, err2 := gh.RESTClient(nil)
	err2 = client.Get(
		"repos/"+repoWithOrg,
		&repo)
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	if config.verbose {
		fmt.Printf("Repository %s has:\n", repoWithOrg)
	} else {
		fmt.Printf("Repo %s has: ", repoWithOrg)
	}

	if len(repo.Description) > 0 {
		if config.verbose {
			fmt.Printf("  - a description ☑️\n")
		} else {
			fmt.Printf("description ☑️, ")
		}
	} else {
		if config.verbose {
			fmt.Printf("  - no description 😇\n")
		} else {
			fmt.Printf("no description 😇, ")
		}
	}

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
		return
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

	// https://docs.github.com/en/rest/reference/collaborators#list-repository-collaborators
	collaborators := []struct {
		login string
	}{}
	client, errCollabs := gh.RESTClient(nil)
	errCollabs = client.Get(
		"repos/"+repoWithOrg+"/collaborators",
		&collaborators)
	if errCollabs != nil && len(errCollabs.Error()) > 0 {
		if strings.HasPrefix(errCollabs.Error(), "HTTP 403") {
			// 🤫
		} else {
			fmt.Println(errReadme)
			return
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

	// https://docs.github.com/en/rest/reference/metrics#get-community-profile-metrics
	if strings.Compare(repo.Visibility, "public") == 0 {
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
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
