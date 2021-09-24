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

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubAccessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	repos, _, err := client.Repositories.List(ctx, "", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	hardcodedRepo := repos[3] // TODO
	owner := *hardcodedRepo.Owner.Login
	repo := *hardcodedRepo.Name
	comp, _, err := client.Repositories.CompareCommits(ctx, owner, repo, "master", "premaster2", nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	var allChangeMessages string
	for _, c := range comp.Commits {
		allChangeMessages += *c.Commit.Message
	}
	//println(allChangeMessages)
	idSHniki := map[string]bool{}
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
			idSHniki[i[1]] = true
		}
		fmt.Printf("%s # %d # %+v\n", pattern, len(idSHniki), idSHniki)
	}

	var resultStroyIds []string
	for stringId := range idSHniki {
		resultStroyIds = append(resultStroyIds, stringId)
	}

	fmt.Printf("# %d # %+v\n", len(resultStroyIds), resultStroyIds)

}
