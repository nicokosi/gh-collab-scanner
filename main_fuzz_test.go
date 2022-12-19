package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func FuzzScanRepo(f *testing.F) {
	defer gock.Off()
	//gock.Observe(gock.DumpRequest)
	defer gock.DisableNetworking()
	gock.New("https://api.github.com").
		Get("/repos/(.*)/readme").
		Reply(200)
	f.Fuzz(func(t *testing.T, randomOrgName string, randomRepoName string) {
		repoWithOrg := randomOrgName + "/" + randomRepoName
		if repoWithOrg == url.QueryEscape(repoWithOrg) {
			gock.New("https://api.github.com").
				Get("/repos/(.*)").
				Reply(200).
				JSON(map[string]string{"name": repoWithOrg})
			message, repository, validRepo := scanRepo(config{repo: repoWithOrg}, repoWithOrg)
			assert.Equal(t, repository.Name, repoWithOrg, "expecting same repo name in output/input")
			assert.True(t, validRepo, "expecting valid repo for "+repoWithOrg)
			assert.NotEmpty(t, message, "expecting non-empty message for "+repoWithOrg)
		}
	})
}
