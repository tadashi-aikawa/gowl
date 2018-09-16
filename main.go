package main

import "fmt"

func main() {
	res, err := SearchRepositories("miroir")
	if err != nil {
		panic(err)
	}

	for _, repo := range res.Items {
		fmt.Printf("★%-8v%-15v%-15v%-45v%v\n", repo.StargazersCount, repo.ID, repo.Language, repo.FullName, repo.UpdatedAt)
	}
}
