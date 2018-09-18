package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

func main() {
	client := github.NewClient(nil)
	repos, _, err := client.Search.Repositories(context.Background(), "miroir", nil)
	if err != nil {
		panic(err)
	}

	for _, repo := range repos.Repositories {
		fmt.Printf("%v\n", repo.GetFullName())
	}

	fmt.Printf("â˜…%-8v%-15v%-15v%-45v%v\n", repo.StargazersCount, repo.ID, repo.Language, repo.FullName, repo.UpdatedAt)
}
