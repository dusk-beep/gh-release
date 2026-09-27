package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gh-release/internal/model"
)

type forgejoVersion struct {
	Version string `json:"version"`
}

func Parse(ctx context.Context, raw string) (model.Repository, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return model.Repository{}, fmt.Errorf("repository URL is empty")
	}

	u, err := url.Parse(raw)
	if err != nil {
		return model.Repository{}, fmt.Errorf("parse repository URL: %w", err)
	}

	if u.Scheme != "https" && u.Scheme != "http" {
		return model.Repository{}, fmt.Errorf("repository URL must use http or https")
	}

	host := strings.ToLower(u.Hostname())
	host = strings.TrimPrefix(host, "www.")

	parts := pathParts(u.Path)
	if len(parts) < 2 {
		return model.Repository{}, fmt.Errorf(
			"invalid repository URL: expected owner/repository",
		)
	}

	switch host {
	case "github.com":
		if len(parts) != 2 {
			return model.Repository{}, fmt.Errorf(
				"invalid GitHub repository URL",
			)
		}

		return model.Repository{
			Provider: "github",
			Host:     host,
			Owner:    parts[0],
			Name:     parts[1],
		}, nil

	case "gitlab.com":
		return model.Repository{
			Provider: "gitlab",
			Host:     host,
			Owner:    strings.Join(parts[:len(parts)-1], "/"),
			Name:     parts[len(parts)-1],
		}, nil

	case "codeberg.org":
		if len(parts) != 2 {
			return model.Repository{}, fmt.Errorf(
				"invalid Codeberg repository URL",
			)
		}

		return model.Repository{
			Provider: "forgejo",
			Host:     host,
			Owner:    parts[0],
			Name:     parts[1],
		}, nil
	}

	// Unknown host: determine whether it is a Forgejo instance.
	if err := checkForgejo(ctx, u); err == nil {
		return model.Repository{
			Provider: "forgejo",
			Host:     host,
			Owner:    strings.Join(parts[:len(parts)-1], "/"),
			Name:     parts[len(parts)-1],
		}, nil
	}

	return model.Repository{}, fmt.Errorf(
		"unsupported repository host %q",
		host,
	)
}

func checkForgejo(ctx context.Context, base *url.URL) error {
	apiURL := *base
	apiURL.Path = "/api/v1/version"
	apiURL.RawPath = ""

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		apiURL.String(),
		nil,
	)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("not a Forgejo API: HTTP %d", resp.StatusCode)
	}

	var version forgejoVersion

	if err := json.NewDecoder(resp.Body).Decode(&version); err != nil {
		return fmt.Errorf("invalid Forgejo version response: %w", err)
	}

	if version.Version == "" {
		return fmt.Errorf("missing Forgejo version")
	}

	return nil
}

func pathParts(path string) []string {
	path = strings.Trim(path, "/")

	parts := strings.Split(path, "/")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		if part == "" {
			continue
		}

		part = strings.TrimSuffix(part, ".git")

		if part != "" {
			result = append(result, part)
		}
	}

	return result
}
