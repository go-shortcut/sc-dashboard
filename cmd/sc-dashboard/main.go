package main

import (
	"context"
	"fmt"
	"github.com/go-shortcut/go-shortcut-api/pkg/shortcutclient"
	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
	"log"
	"os"
	"regexp"
	"strconv"
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
	token := os.Getenv("SHORTCUT_API_TOKEN")
	if token == "" {
		fmt.Println("SHORTCUT_API_TOKEN environ is required")
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

	storyIds := GetKeysAsInt64Slice(scStoryIds)
	fmt.Printf("# %d # %d # %+v\n", len(scStoryIds), len(storyIds), storyIds)

	GithubBranchName := "premaster2"

	shortcutClient := shortcutclient.New(token)
	//shortcutClient.HTTPClient.Timeout = 30 * time.Second
	//shortcutClient.Debug = true
	playload := shortcutclient.UpdateMultipleStoriesParams{
		StoryIds: storyIds,
		LabelsAdd: []shortcutclient.CreateLabelParams{
			{Name: GithubBranchName, Color: "#0000FF"},
		},
	}

	if GithubBranchName == "master" {
		playload.LabelsRemove = []shortcutclient.CreateLabelParams{
			{Name: "premaster2", Color: "#FF0000"},
		}
	}

	if _, err = shortcutClient.UpdateMultipleStories(playload); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

}

func GetKeysAsInt64Slice(m map[string]interface{}) []int64 {
	keys := make([]int64, len(m))

	i := 0
	for k := range m {
		kInt64, e := strconv.ParseInt(k, 10, 64)
		if e != nil {
			continue
		}
		keys[i] = kInt64
		i++
	}
	return keys
}
