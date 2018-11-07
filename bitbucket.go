package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// BitbucketClient includes fields used by gowl
type BitbucketClient struct {
	BaseURL  string
	UserName string
	Password string
}

// BitbucketRepositoryResponse is struct of a API response
type BitbucketRepositoryResponse struct {
	Size       int64                 `json:"size"`
	Limit      int64                 `json:"limit"`
	IsLastPage bool                  `json:"isLastPage"`
	Values     []BitbucketRepository `json:"values"`
	Start      int64                 `json:"start"`
}

// BitbucketRepository includes fields used by gowl
type BitbucketRepository struct {
	Slug          string     `json:"slug"`
	ID            int64      `json:"id"`
	Name          string     `json:"name"`
	SCMID         string     `json:"scmId"`
	State         string     `json:"state"`
	StatusMessage string     `json:"statusMessage"`
	Forkable      bool       `json:"forkable"`
	Project       project    `json:"project"`
	Public        bool       `json:"public"`
	Links         valueLinks `json:"links"`
}

// GetFullName returns project/repository
func (r *BitbucketRepository) GetFullName() string {
	// Slug is always lowercase
	return fmt.Sprintf("%v/%v", strings.ToLower(r.Project.Key), r.Slug)
}

type valueLinks struct {
	Clone []cloneSeed `json:"clone"`
	Self  []self      `json:"self"`
}

type cloneSeed struct {
	Href string `json:"href"`
	Name string `json:"name"`
}

type self struct {
	Href string `json:"href"`
}

type project struct {
	Key         string       `json:"key"`
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Public      bool         `json:"public"`
	Type        string       `json:"type"`
	Links       projectLinks `json:"links"`
}

type projectLinks struct {
	Self []self `json:"self"`
}

func (c *BitbucketClient) searchRepositories(word string) (BitbucketRepositoryResponse, error) {
	url := fmt.Sprintf("%v/rest/api/1.0/repos?name=%v", c.BaseURL, word)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return BitbucketRepositoryResponse{}, err
	}

	req.SetBasicAuth(c.UserName, c.Password)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return BitbucketRepositoryResponse{}, err
	}
	defer res.Body.Close()

	var r BitbucketRepositoryResponse
	json.NewDecoder(res.Body).Decode(&r)

	return r, nil
}
