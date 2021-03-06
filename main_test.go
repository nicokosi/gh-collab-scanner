package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestScanRepo(t *testing.T) {
	defer gock.Off()
	defer gock.DisableNetworking()
	gock.New("https://api.github.com").
		Get("/repos/acme/buzz/readme").
		Reply(200).
		JSON(map[string]string{"name": "buzz"})
	gock.New("https://api.github.com").
		Get("/repos/acme/buzz").
		Reply(200).
		File("test_repo.json")

	message, repository, validRepo := scanRepo(config{repo: "acme/buzz"}, "acme/buzz")

	assert.True(t, validRepo)
	assert.Equal(t, repo{Name: "buzz", Owner: owner{Login: "Coyote"}, Description: "Beep, beep", Topics: []string{"old", "cartoon"}, Visibility: "public"}, repository)
	assert.Equal(t, "acme/buzz: README ☑️, description ☑️, topics ☑️, ", message)
}

func TestGetRepos_for_org(t *testing.T) {
	defer gock.Off()
	defer gock.DisableNetworking()
	gock.New("https://api.github.com").
		Get("/orgs/acme/repos").
		Reply(200).
		File("test_repos.json")

	repositories, error := getRepos(config{org: "acme"})

	assert.Len(t, repositories, 2)
	assert.Nil(t, error)
}

func TestGetRepos_for_org_with_topic(t *testing.T) {
	defer gock.Off()
	defer gock.DisableNetworking()
	gock.New("https://api.github.com").
		Get("/orgs/acme/repos").
		Reply(200).
		File("test_repos.json")

	repositories, error := getRepos(config{org: "acme", topic: "rabbit"})

	assert.Len(t, repositories, 1)
	assert.Nil(t, error)
}

func TestGetRepos_for_org_with_missing_topic(t *testing.T) {
	defer gock.Off()
	defer gock.DisableNetworking()
	gock.New("https://api.github.com").
		Get("/orgs/acme/repos").
		Reply(200).
		File("test_repos.json")

	repositories, error := getRepos(config{org: "acme", topic: "does-not-exist"})

	assert.Empty(t, repositories)
	assert.Nil(t, error)
}

func TestGetRepos_for_user(t *testing.T) {
	defer gock.Off()
	defer gock.DisableNetworking()
	gock.New("https://api.github.com").
		Get("/users/acme/repos").
		Reply(200).
		File("test_repos.json")

	repositories, error := getRepos(config{user: "acme"})

	assert.Len(t, repositories, 2)
	assert.Nil(t, error)
}

func TestGetRepos_for_user_with_topic(t *testing.T) {
	defer gock.Off()
	defer gock.DisableNetworking()
	gock.New("https://api.github.com").
		Get("/users/acme/repos").
		Reply(200).
		File("test_repos.json")

	repositories, error := getRepos(config{user: "acme", topic: "rabbit"})

	assert.Len(t, repositories, 1)
	assert.Nil(t, error)
}

func TestGetRepos_for_user_with_missing_topic(t *testing.T) {
	defer gock.Off()
	defer gock.DisableNetworking()
	gock.New("https://api.github.com").
		Get("/users/acme/repos").
		Reply(200).
		File("test_repos.json")

	repositories, error := getRepos(config{user: "acme", topic: "does-not-exist"})

	assert.Empty(t, repositories)
	assert.Nil(t, error)
}

func TestGetRepos_error(t *testing.T) {
	defer gock.Off()
	defer gock.DisableNetworking()
	gock.New("https://api.github.com").
		Get("/users/acme/repos").
		Reply(500)

	repositories, error := getRepos(config{user: "acme"})

	assert.Empty(t, repositories)
	assert.NotNil(t, error)
}

func TestGetRepos_no_user_nor_org(t *testing.T) {
	defer gock.Off()

	repositories, error := getRepos(config{})

	assert.Empty(t, repositories)
	assert.Nil(t, error)
}

func TestScanRepo_Org(t *testing.T) {
	defer gock.Off()
	defer gock.DisableNetworking()
	gock.New("https://api.github.com").
		Get("/repos/acme/buzz/readme").
		Reply(200).
		JSON(map[string]string{"name": "buzz"})
	gock.New("https://api.github.com").
		Get("/repos/acme/buzz").
		Reply(200).
		File("test_repo.json")

	message, repository, validRepo := scanRepo(config{org: "acme"}, "acme/buzz")

	assert.True(t, validRepo)
	assert.Equal(t, repo{Name: "buzz", Owner: owner{Login: "Coyote"}, Description: "Beep, beep", Topics: []string{"old", "cartoon"}, Visibility: "public"}, repository)
	assert.Equal(t, "acme/buzz: README ☑️, description ☑️, topics ☑️, ", message)
}

func TestScanRepo_Verbose(t *testing.T) {
	defer gock.Off()
	defer gock.DisableNetworking()
	gock.New("https://api.github.com").
		Get("/repos/acme/buzz/readme").
		Reply(200).
		JSON(map[string]string{"name": "buzz"})
	gock.New("https://api.github.com").
		Get("/repos/acme/buzz").
		Reply(200).
		File("test_repo.json")

	message, repository, validRepo := scanRepo(config{repo: "acme/buzz", verbose: true}, "acme/buzz")

	assert.True(t, validRepo)
	assert.Equal(t, repo{Name: "buzz", Owner: owner{Login: "Coyote"}, Description: "Beep, beep", Topics: []string{"old", "cartoon"}, Visibility: "public"}, repository)
	assert.Equal(t, "acme/buzz has: \n  - a README ☑️\n  - a description ☑️\n  - topics ☑️", message)
}

func TestScanRepo_Verbose_ReadmeError(t *testing.T) {
	defer gock.Off()
	defer gock.DisableNetworking()
	gock.New("https://api.github.com").
		Get("/repos/acme/buzz/readme").
		Reply(500)
	gock.New("https://api.github.com").
		Get("/repos/acme/buzz").
		Reply(200).
		File("test_repo.json")

	message, repository, validRepo := scanRepo(config{repo: "acme/buzz", verbose: true}, "acme/buzz")

	assert.True(t, validRepo)
	assert.Equal(t, repo{Name: "buzz", Owner: owner{Login: "Coyote"}, Description: "Beep, beep", Topics: []string{"old", "cartoon"}, Visibility: "public"}, repository)
	assert.Equal(t, "acme/buzz has: \n  - a description ☑️\n  - topics ☑️", message)
}

func TestScanRepo_Verbose_RepoError(t *testing.T) {
	defer gock.Off()
	defer gock.DisableNetworking()
	gock.New("https://api.github.com").
		Get("/repos/acme/buzz/readme").
		Reply(200).
		JSON(map[string]string{"name": "buzz"})
	gock.New("https://api.github.com").
		Get("/repos/acme/buzz").
		Reply(500)

	message, _, validRepo := scanRepo(config{repo: "acme/buzz", verbose: true}, "acme/buzz")

	assert.False(t, validRepo)
	assert.Equal(t, "acme/buzz has: \n  - a README ☑️", message)
}

func TestScanCollaborators(t *testing.T) {
	defer gock.Off()
	defer gock.DisableNetworking()
	gock.New("https://api.github.com").
		Get("/repos/acme/buzz/collaborators").
		Reply(200).
		File("test_collaborators.json")

	message := scanCollaborators(config{repo: "acme/buzz"}, "acme/buzz")

	assert.Equal(t, "1 collaborator 👤, ", message)
}

func TestScanCollaborators_Verbose(t *testing.T) {
	defer gock.Off()
	defer gock.DisableNetworking()
	gock.New("https://api.github.com").
		Get("/repos/acme/buzz/collaborators").
		Reply(200).
		File("test_collaborators.json")

	message := scanCollaborators(config{repo: "acme/buzz", verbose: true}, "acme/buzz")

	assert.Equal(t, "\n  - 1 collaborator 👤", message)
}

func TestScanCommunityScore(t *testing.T) {
	defer gock.Off()
	defer gock.DisableNetworking()
	gock.New("https://api.github.com").
		Get("/repos/acme/buzz/community/profile").
		Reply(200).
		JSON(map[string]int{"health_percentage": 42})

	message := scanCommunityScore(config{repo: "acme/buzz"}, "acme/buzz")

	assert.Equal(t, "community profile score: 42 💯", message)
}

func TestScanCommunityScore_Verbose(t *testing.T) {
	defer gock.Off()
	defer gock.DisableNetworking()
	gock.New("https://api.github.com").
		Get("/repos/acme/buzz/community/profile").
		Reply(200).
		JSON(map[string]int{"health_percentage": 42})

	message := scanCommunityScore(config{repo: "acme/buzz", verbose: true}, "acme/buzz")

	assert.Equal(t, "\n  - a community profile score of 42 💯", message)
}

func TestVersion(t *testing.T) {
	version := getVersion()

	assert.NotEmpty(t, version.commit)
	assert.NotEmpty(t, version.date)
}
