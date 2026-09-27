package provider

import (
	"context"
	"fmt"

	"gh-release/internal/model"
)

type Provider interface {
	Name() string
	LatestRelease(ctx context.Context, repo model.Repository) (model.Release, error)
}

type Registry struct {
	providers map[string]Provider
}

func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

func (r *Registry) Register(p Provider) {
	r.providers[p.Name()] = p
}

func (r *Registry) Get(name string) (Provider, error) {
	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %q not registered", name)
	}

	return p, nil
}
