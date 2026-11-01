package model

import "time"

type Repository struct {
	Provider string
	Host     string
	Owner    string
	Name     string
}

func (r Repository) FullName() string {
	return r.Owner + "/" + r.Name
}

func (r Repository) Key() string {
	host := r.Host
	if host == "" {
		host = r.Provider
	}

	return host + "/" + r.Owner + "/" + r.Name
}

func (r Repository) URL() string {
	host := r.Host
	if host == "" {
		switch r.Provider {
		case "github":
			host = "github.com"
		case "gitlab":
			host = "gitlab.com"
		default:
			host = r.Provider
		}
	}

	return "https://" + host + "/" + r.Owner + "/" + r.Name
}

type Release struct {
	ID          string
	Repository  Repository
	Tag         string
	Name        string
	Body        string
	URL         string
	PublishedAt time.Time
	CreatedAt   time.Time
	Prerelease  bool
	Draft       bool
}

type RepositoryResult struct {
	Repository Repository
	Release    Release
	Err        error
}
