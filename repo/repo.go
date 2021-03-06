package repo

import (
	"time"
)

type Client interface {
	GetPublicRepos() ([]*Repo, error)
	GetContributions() ([]*Repo, error)
}

type Repo struct {
	Name, URL, Repo, Source string
	Contribution            bool
}

type Config struct {
	Github RepoConfig
	Gitlab RepoConfig
}

type RepoConfig struct {
	Username string
	Token    string
}

type RepoClient struct {
	clients       []Client
	repos         []*Repo
	contributions []*Repo
}

func GetClient(c *Config) *RepoClient {
	return &RepoClient{
		repos:         make([]*Repo, 0),
		contributions: make([]*Repo, 0),
		clients: []Client{
			getGitHubClient(c),
			getGitlabClient(c),
		},
	}
}

func (c *RepoClient) GetRepos(fresh bool) ([]*Repo, error) {
	if !fresh {
		return c.repos, nil
	}

	allRepos := make([]*Repo, 0)
	for _, c := range c.clients {
		repos, err := c.GetPublicRepos()
		if err != nil {
			return nil, err
		}

		if len(repos) > 0 {
			allRepos = append(allRepos, repos...)
		}

		repos, err = c.GetContributions()
		if err != nil {
			return nil, err
		}

		if len(repos) > 0 {
			allRepos = append(allRepos, repos...)
		}
	}

	c.repos = allRepos
	return c.repos, nil
}

func (c *RepoClient) Poll(d time.Duration) {
	c.GetRepos(true)
	ticker := time.NewTicker(d)
	for range ticker.C {
		c.GetRepos(true)
	}
}
