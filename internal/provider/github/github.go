package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v91/github"

	"gh-release/internal/model"
	"gh-release/internal/provider"
)

type Provider struct {
	client *github.Client
}

func New(token string) (*Provider, error) {
	var options []github.ClientOptionsFunc

	if token != "" {
		options = append(options, github.WithAuthToken(token))
	}

	client, err := github.NewClient(options...)
	if err != nil {
		return nil, fmt.Errorf("create github client: %w", err)
	}

	return &Provider{client: client}, nil
}

func (p *Provider) Name() string {
	return "github"
}

func (p *Provider) LatestRelease(
	ctx context.Context,
	repo model.Repository,
) (model.Release, error) {
	r, _, err := p.client.Repositories.GetLatestRelease(
		ctx,
		repo.Owner,
		repo.Name,
	)
	if err != nil {
		return model.Release{}, fmt.Errorf("get latest GitHub release: %w", err)
	}

	return model.Release{
		ID:          fmt.Sprintf("%d", r.ID),
		Repository:  repo,
		Tag:         r.TagName,
		Name:        stringValue(r.Name),
		Body:        *r.Body,
		URL:         r.HTMLURL,
		PublishedAt: timestampValue(r.PublishedAt),
		CreatedAt:   r.CreatedAt.Time,
		Prerelease:  r.Prerelease,
		Draft:       r.Draft,
	}, nil
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func timestampValue(value *github.Timestamp) time.Time {
	if value == nil {
		return time.Time{}
	}

	return value.Time
}

var _ provider.Provider = (*Provider)(nil)
