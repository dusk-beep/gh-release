package gitlab

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"gh-release/internal/model"
	"gh-release/internal/provider"
)

type Provider struct {
	client *http.Client
	token  string
}

func New(token string) *Provider {
	return &Provider{
		client: http.DefaultClient,
		token:  token,
	}
}

func (p *Provider) Name() string {
	return "gitlab"
}

type releaseResponse struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	CreatedAt   string `json:"created_at"`
	Description string `json:"description"`
	ReleasedAt  string `json:"released_at"`
	Links       struct {
		Web string `json:"web"`
	} `json:"_links"`
}

func (p *Provider) LatestRelease(
	ctx context.Context,
	repo model.Repository,
) (model.Release, error) {
	project := url.PathEscape(repo.Owner + "/" + repo.Name)

	endpoint := fmt.Sprintf(
		"https://gitlab.com/api/v4/projects/%s/releases/permalink/latest",
		project,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return model.Release{}, fmt.Errorf("create GitLab request: %w", err)
	}

	if p.token != "" {
		req.Header.Set("PRIVATE-TOKEN", p.token)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return model.Release{}, fmt.Errorf("request GitLab release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Release{}, fmt.Errorf(
			"GitLab returned HTTP %d",
			resp.StatusCode,
		)
	}

	var release releaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return model.Release{}, fmt.Errorf("decode GitLab release: %w", err)
	}

	publishedAt, err := parseTime(release.ReleasedAt)
	if err != nil {
		return model.Release{}, fmt.Errorf("parse GitLab release time: %w", err)
	}

	createdAt, err := parseTime(release.CreatedAt)
	if err != nil {
		return model.Release{}, fmt.Errorf("parse GitLab creation time: %w", err)
	}

	return model.Release{
		ID:          release.TagName,
		Repository:  repo,
		Body:        release.Description,
		Tag:         release.TagName,
		Name:        release.Name,
		URL:         release.Links.Web,
		PublishedAt: publishedAt,
		CreatedAt:   createdAt,
	}, nil
}

func parseTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}

	return time.Parse(time.RFC3339, value)
}

var _ provider.Provider = (*Provider)(nil)
