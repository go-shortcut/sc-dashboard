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
	"strings"
)

const (
	EnvKeyGithubBaseRef = "GITHUB_BASE_REF"
	EnvKeyGithubHeadRef = "GITHUB_HEAD_REF"
)

func main() {
	githubAccessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	if githubAccessToken == "" {
		fmt.Println("GITHUB_ACCESS_TOKEN environ is required")
		os.Exit(1)
	}

	token := os.Getenv("SHORTCUT_API_TOKEN")
	if token == "" {
		fmt.Println("SHORTCUT_API_TOKEN environ is required")
		os.Exit(1)
	}

	GithubRepository := os.Getenv("GITHUB_REPOSITORY")
	var githubOwnerName, githubRepoName string
	if GithubRepository == "" {
		fmt.Println("GITHUB_REPOSITORY environ is required")
		fmt.Println("GITHUB_REPOSITORY\tThe owner and repository name. For example, octocat/Hello-World.")
		os.Exit(1)
	} else {
		splitedGithubRepository := strings.Split(GithubRepository, "/")
		if len(splitedGithubRepository) != 2 {
			fmt.Println("Failed split GITHUB_REPOSITORY to the owner and repository name.")
			os.Exit(1)
		}
		githubOwnerName, githubRepoName = splitedGithubRepository[0], splitedGithubRepository[1]
	}

	GithubBaseRef := os.Getenv(EnvKeyGithubBaseRef)
	GithubHeadRef := os.Getenv(EnvKeyGithubHeadRef)
	if len(GithubBaseRef)*len(GithubHeadRef) == 0 {
		fmt.Printf("Bad ref variables: %s=%s, %s=%s.\n", EnvKeyGithubBaseRef, GithubBaseRef, EnvKeyGithubHeadRef, GithubHeadRef)
		os.Exit(1)

	}

	ShortcutAddLabel := os.Getenv("SHORTCUT_ADD_LABEL")
	ShortcutDelLabel := os.Getenv("SHORTCUT_DEL_LABEL")
	if len(ShortcutAddLabel)+len(ShortcutDelLabel) == 0 {
		// SHORTCUT_ADD_LABEL="master marty" SHORTCUT_DEL_LABEL="premaster2 marty"
		fmt.Println("At least one variable must be set: SHORTCUT_ADD_LABEL, SHORTCUT_DEL_LABEL.")
		os.Exit(1)

	}

	// all env variables are checked, let's go
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubAccessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	comp, _, err := client.Repositories.CompareCommits(ctx, githubOwnerName, githubRepoName, GithubBaseRef, GithubHeadRef, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	if len(comp.Commits) == 0 {
		fmt.Printf("No commits found between %s and %s.\n", GithubBaseRef, GithubHeadRef)
		os.Exit(0)
	}
	fmt.Printf("Found %+v commits.\n", len(comp.Commits))

	var allChangeMessages string
	for _, c := range comp.Commits {
		allChangeMessages += *(c.Commit.Message)
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
	if len(scStoryIds) == 0 {
		fmt.Println("No story ids found")
		os.Exit(0)
	}

	storyIds := GetKeysAsInt64Slice(scStoryIds)
	fmt.Printf("Found story ids: %+v", scStoryIds)

	shortcutClient := shortcutclient.New(token)
	//shortcutClient.HTTPClient.Timeout = 30 * time.Second
	//shortcutClient.Debug = true
	playload := shortcutclient.UpdateMultipleStoriesParams{
		StoryIds: storyIds,
	}

	if ShortcutAddLabel != "" {
		playload.LabelsAdd = []shortcutclient.CreateLabelParams{
			{Name: ShortcutAddLabel, Color: "#0000FF"},
		}
	}
	if ShortcutDelLabel != "" {
		playload.LabelsRemove = []shortcutclient.CreateLabelParams{
			{Name: ShortcutDelLabel},
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
