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
	"time"
)

const (
	// Default environment variables
	// https://docs.github.com/en/actions/learn-github-actions/environment-variables
	envKeyGithubRepository = "GITHUB_REPOSITORY"

	// Secrets
	envKeyGithubAccessToken = "GITHUB_ACCESS_TOKEN"
	envKeyShortcutApiToken  = "SHORTCUT_API_TOKEN"

	// Pull request number from $GITHUB_EVENT_PATH file
	// or template variable ${{ github.event.pull_request.number }}
	envKeyPullNumber = "PULL_NUMBER"

	// user defined
	envKeyShortcutAddLabel = "SHORTCUT_ADD_LABEL"
	envKeyShortcutDelLabel = "SHORTCUT_DEL_LABEL"
)

func main() {
	githubAccessToken := os.Getenv(envKeyGithubAccessToken)
	if githubAccessToken == "" {
		fmt.Println(envKeyGithubAccessToken + " environ is required")
		os.Exit(1)
	}

	shortcutApiToken := os.Getenv(envKeyShortcutApiToken)
	if shortcutApiToken == "" {
		fmt.Println(envKeyShortcutApiToken + " environ is required")
		os.Exit(1)
	}

	githubRepository := os.Getenv(envKeyGithubRepository)
	var githubOwnerName, githubRepoName string
	if githubRepository == "" {
		fmt.Println(envKeyGithubRepository + " environ is required")
		fmt.Println(envKeyGithubRepository + "\tThe owner and repository name. For example, octocat/Hello-World.")
		os.Exit(1)
	} else {
		splitGithubRepository := strings.Split(githubRepository, "/")
		if len(splitGithubRepository) != 2 {
			fmt.Println("Failed split " + envKeyGithubRepository + " to the owner and repository name.")
			os.Exit(1)
		}
		githubOwnerName, githubRepoName = splitGithubRepository[0], splitGithubRepository[1]
	}

	shortcutAddLabel := os.Getenv(envKeyShortcutAddLabel)
	shortcutDelLabel := os.Getenv(envKeyShortcutDelLabel)
	if len(shortcutAddLabel)+len(shortcutDelLabel) == 0 {
		fmt.Println("At least one variable must be set: " + envKeyShortcutAddLabel + "," + envKeyShortcutDelLabel)
		os.Exit(1)

	}

	githubPullNumber := os.Getenv(envKeyPullNumber)
	pNum, e := strconv.Atoi(githubPullNumber)
	if e != nil {
		fmt.Printf("%s=%s is not int64.\n", envKeyPullNumber, githubPullNumber)
		os.Exit(1)
	}

	// all env variables are checked, let's go
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubAccessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	//browserOpts := github.ListOptions{PerPage: 100}
	curPR, _, err := client.PullRequests.Get(ctx, githubOwnerName, githubRepoName, pNum)
	if err != nil {
		fmt.Printf("Pull request not found: %d.\n", pNum)
		os.Exit(1)
	}

	GithubBaseRef := *(curPR.Base.Ref)
	GithubHeadRef := *(curPR.Head.Ref)
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
	fmt.Printf("Found story ids: %+v.\n", storyIds)

	shortcutClient := shortcutclient.New(shortcutApiToken)
	shortcutClient.HTTPClient.Timeout = 30 * time.Second
	//shortcutClient.Debug = true
	playload := shortcutclient.UpdateMultipleStoriesParams{
		StoryIds: storyIds,
	}

	if shortcutAddLabel != "" {
		playload.LabelsAdd = []shortcutclient.CreateLabelParams{
			{Name: shortcutAddLabel, Color: "#0000FF"},
		}
	}
	if shortcutDelLabel != "" {
		playload.LabelsRemove = []shortcutclient.CreateLabelParams{
			{Name: shortcutDelLabel},
		}
	}

	updatedStroies, err := shortcutClient.UpdateMultipleStories(playload)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	countStories := len(updatedStroies)
	storyIDs := make([]int64, countStories)
	messages := make(map[int64]string, countStories)

	for i, s := range updatedStroies {
		storyIDs[i] = s.ID
		messages[s.EpicID] += " - [" + s.Name + "](" + s.AppURL + ")\n"

	}
	fmt.Printf("updated %v stories: %v.\n", countStories, storyIDs)

	scEpics, err := shortcutClient.ListEpics()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	scEpicTitle := make(map[int64]string, len(scEpics))
	for _, scEpic := range scEpics {
		scEpicTitle[scEpic.ID] = "# [" + scEpic.Name + "](" + scEpic.AppURL + ")\n"
	}
	//fmt.Printf("%+v", scEpicTitle)
	var body string

	for epicId, urls := range messages {
		body += scEpicTitle[epicId] + "\n" + urls + "\n\n"
	}

	_, _, err = client.Issues.CreateComment(ctx, githubOwnerName, githubRepoName, pNum,
		&github.IssueComment{Body: &body},
	)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(body)
}
