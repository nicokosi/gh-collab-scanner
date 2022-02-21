package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/cli/go-gh"
)

func main() {

	var repoWithOrg = ""
	if len(os.Args) > 1 {
		repoWithOrg = os.Args[1]
	} else {
		fmt.Printf("(current repo)\n")
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
	if err != nil {
		fmt.Println(err)
		return
	}

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
	fmt.Printf("Repository %s has:\n", repoWithOrg)
	if len(repo.Description) > 0 {
		fmt.Printf("  - a description ☑️\n")
	} else {
		fmt.Printf("  - no description 😇\n")
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
		fmt.Printf("  - has a README ☑️\n")
	} else if strings.HasPrefix(errReadme.Error(), "HTTP 404: Not Found") {
		fmt.Printf("  - no README 😇\n")
	} else {
		fmt.Println(errReadme)
		return
	}

	if len(repo.Topics) > 0 {
		fmt.Printf("  - topics ☑️\n")
	} else {
		fmt.Printf("  - no topics 😇\n")
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
		fmt.Printf("  - %d collaborator 👤\n", len(collaborators))
	} else {
		fmt.Printf("  - %d collaborators 👥\n", len(collaborators))
	}

	// https://docs.github.com/en/rest/reference/metrics#get-community-profile-metrics
	if strings.Compare(repo.Visibility, "public") == 0 {
		community_profile := struct {
			Health_percentage int64
		}{}
		client, errCommunityProfile := gh.RESTClient(nil)
		errCommunityProfile = client.Get(
			"repos/"+repoWithOrg+"/community/profile",
			&community_profile)
		if errCommunityProfile != nil {
			fmt.Println(errCommunityProfile)
		}
		fmt.Printf("  - a community profile score of %d 💯\n", community_profile.Health_percentage)
	}
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
