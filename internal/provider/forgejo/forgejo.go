package forgejo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
	return "forgejo"
}

type releaseResponse struct {
	ID          int64      `json:"id"`
	TagName     string     `json:"tag_name"`
	Name        string     `json:"name"`
	Body        string     `json:"body"`
	HTMLURL     string     `json:"html_url"`
	CreatedAt   *time.Time `json:"created_at"`
	PublishedAt *time.Time `json:"published_at"`
	Prerelease  bool       `json:"prerelease"`
	Draft       bool       `json:"draft"`
}

func (p *Provider) LatestRelease(
	ctx context.Context,
	repo model.Repository,
) (model.Release, error) {
	endpoint := fmt.Sprintf(
		"https://%s/api/v1/repos/%s/%s/releases/latest",
		repo.Host,
		repo.Owner,
		repo.Name,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return model.Release{}, fmt.Errorf("create Forgejo request: %w", err)
	}

	if p.token != "" {
		req.Header.Set("Authorization", "token "+p.token)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return model.Release{}, fmt.Errorf("request Forgejo release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Release{}, fmt.Errorf(
			"Forgejo returned HTTP %d",
			resp.StatusCode,
		)
	}

	var release releaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return model.Release{}, fmt.Errorf("decode Forgejo release: %w", err)
	}

	result := model.Release{
		ID:         fmt.Sprintf("%d", release.ID),
		Repository: repo,
		Tag:        release.TagName,
		Body:       release.Body,
		Name:       release.Name,
		URL:        release.HTMLURL,
		Prerelease: release.Prerelease,
		Draft:      release.Draft,
	}

	if release.CreatedAt != nil {
		result.CreatedAt = release.CreatedAt.UTC()
	}

	if release.PublishedAt != nil {
		result.PublishedAt = release.PublishedAt.UTC()
	}

	return result, nil
}

var _ provider.Provider = (*Provider)(nil)
