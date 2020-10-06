package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/schigh/tools/pkg/github"
)

type Schema struct {
	Data Data `json:"data"`
}
type ReleaseNode struct {
	TagName      string `json:"tagName"`
	IsPrerelease bool   `json:"isPrerelease"`
}
type Releases struct {
	Nodes []ReleaseNode `json:"nodes"`
}
type RepoNode struct {
	Name     string   `json:"name"`
	Releases Releases `json:"releases"`
}
type Repositories struct {
	Nodes []RepoNode `json:"nodes"`
}
type Team struct {
	Description  string       `json:"description"`
	Repositories Repositories `json:"repositories"`
}
type Organization struct {
	Team Team `json:"team"`
}
type Data struct {
	Organization Organization `json:"organization"`
}

var (
	team  string
	token string
	org   string
)

const payload = `
{
	"query":"{organization(login: \"%s\") {\n team(slug: \"%s\") {\n description\n repositories(first: 100, orderBy: {field: NAME, direction: DESC}) {\n nodes {\n name\n releases(first: 1, orderBy: {field: CREATED_AT, direction: DESC}) {\n nodes {\n tagName\n isPrerelease \n} \n} \n} \n} \n} \n} \n}"
}
`
const endpoint = "https://api.github.com/graphql"

func main() {
	flag.StringVar(&token, "token", os.Getenv("GITHUB_API_TOKEN"), "github api token")
	flag.StringVar(&team, "team", os.Getenv("GITHUB_TEAM_ID"), "github team ID")
	flag.StringVar(&org, "org", os.Getenv("GITHUB_ORG_ID"), "github org login")
	flag.Parse()

	data := []byte(fmt.Sprintf(payload, org, team))

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "bearer "+token)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, respErr := client.Do(req)
	if respErr != nil {
		log.Fatal(respErr)
	}
	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	schema := &github.Schema{}
	if jsonErr := json.Unmarshal(body, schema); jsonErr != nil {
		log.Fatal(jsonErr)
	}

	repos := schema.Data.Organization.Team.Repositories.Nodes
	fmt.Printf("Showing releases for %s/%s:\n", org, team)
	for _, repo := range repos {
		if len(repo.Releases.Nodes) == 0 {
			continue
		}
		release := repo.Releases.Nodes[0]
		fmt.Printf("%-45s%5s\n", repo.Name, release.TagName)
	}
}
