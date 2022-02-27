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
	assert.Equal(t, "README ☑️, topics ☑️, ", message)
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
	assert.Equal(t, "  - a README ☑️\n  - topics ☑️\n", message)
}

func TestScanRepo_ReadmeError(t *testing.T) {
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
	assert.Equal(t, "  - topics ☑️\n", message)
}

func TestScanRepo_RepoError(t *testing.T) {
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
	assert.Equal(t, "  - a README ☑️\n", message)
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

	assert.Equal(t, "  - 1 collaborator 👤\n", message)
}

func TestScanCommunityScore(t *testing.T) {
	defer gock.Off()
	defer gock.DisableNetworking()
	gock.New("https://api.github.com").
		Get("/repos/acme/buzz/community/profile").
		Reply(200).
		JSON(map[string]int{"health_percentage": 42})

	message := scanCommunityScore(config{repo: "acme/buzz"}, "acme/buzz")

	assert.Equal(t, "community profile score: 42 💯\n", message)
}

func TestScanCommunityScore_Verbose(t *testing.T) {
	defer gock.Off()
	defer gock.DisableNetworking()
	gock.New("https://api.github.com").
		Get("/repos/acme/buzz/community/profile").
		Reply(200).
		JSON(map[string]int{"health_percentage": 42})

	message := scanCommunityScore(config{repo: "acme/buzz", verbose: true}, "acme/buzz")

	assert.Equal(t, "  - a community profile score of 42 💯\n", message)
}