package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
	"log"
	"os"
	"regexp"
)

func main() {
	githubAccessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	if githubAccessToken == "" {
		fmt.Println("GITHUB_ACCESS_TOKEN environ is required")
		os.Exit(1)
	}
	githubRepoName := os.Getenv("GITHUB_REPO_NAME")
	if githubAccessToken == "" {
		fmt.Println("GITHUB_REPO_NAME environ is required")
		os.Exit(1)
	}
	githubOwnerName := os.Getenv("GITHUB_OWNER_NAME")
	if githubAccessToken == "" {
		fmt.Println("GITHUB_OWNER_NAME environ is required")
		os.Exit(1)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubAccessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	comp, _, err := client.Repositories.CompareCommits(ctx, githubOwnerName, githubRepoName, "master", "premaster2", nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	var allChangeMessages string
	for _, c := range comp.Commits {
		allChangeMessages += *c.Commit.Message
	}

	scStoryIds := make(map[string]interface{})
	var chPatterns = []string{
		`/ch([0-9]+)/`,
		`/sc-([0-9]+)/`,
		`\[ch([0-9]+)\]`,
		`\[sc-([0-9]+)\]`,
		`/story/([0-9]+)/`,
	}
	for _, pattern := range chPatterns {
		reCHinBracket := regexp.MustCompile(pattern)
		listReCHinBracket := reCHinBracket.FindAllStringSubmatch(allChangeMessages, -1)
		for _, i := range listReCHinBracket {
			scStoryIds[i[1]] = nil
		}
	}

	fmt.Printf("# %d # %+v\n", len(scStoryIds), GetSliceOfStringKeysFromMap(scStoryIds))

}

func GetSliceOfStringKeysFromMap(m map[string]interface{}) []string {
	keys := make([]string, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
