// it gets the list of branches
package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v49/github"
	"golang.org/x/oauth2"
	"log"
	"os"
	"strings"
)

const (
	envKeyGithubRepository  = "GITHUB_REPOSITORY"
	envKeyGithubAccessToken = "GITHUB_ACCESS_TOKEN"
)

func main() {
	githubAccessToken := os.Getenv(envKeyGithubAccessToken)
	if githubAccessToken == "" {
		fmt.Println(envKeyGithubAccessToken + " environ is required")
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

	// all env variables are checked, let's go
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubAccessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	var allthebranches []*github.Branch
	branchesPage := 1
	protectedBranches := false
	for {
		brz, _, err := client.Repositories.ListBranches(ctx, githubOwnerName, githubRepoName, &github.BranchListOptions{
			Protected:   &protectedBranches,
			ListOptions: github.ListOptions{Page: branchesPage, PerPage: 100},
		})
		if err != nil {
			log.Fatalln(err)
		}
		if len(brz) == 0 {
			break
		}
		allthebranches = append(allthebranches, brz...)
		branchesPage += 1
		log.Println(branchesPage)
	}

	for i, br := range allthebranches {
		fmt.Printf("%v\t%s\t%v\n", i, br.Commit.GetURL(), *br.Name)
	}
}
