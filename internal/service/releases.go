package service

import (
	"context"

	"gh-release/internal/model"
	"gh-release/internal/provider"

	"golang.org/x/sync/errgroup"
)

type Service struct {
	registry    *provider.Registry
	concurrency int
}

func New(registry *provider.Registry, concurrency int) *Service {
	if concurrency < 1 {
		concurrency = 1
	}

	return &Service{
		registry:    registry,
		concurrency: concurrency,
	}
}

func (s *Service) Fetch(
	ctx context.Context,
	repos []model.Repository,
) []model.RepositoryResult {
	results := make([]model.RepositoryResult, len(repos))

	g, ctx := errgroup.WithContext(ctx)
	sem := make(chan struct{}, s.concurrency)

	for i, repo := range repos {
		i, repo := i, repo

		g.Go(func() error {
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				results[i] = model.RepositoryResult{
					Repository: repo,
					Err:        ctx.Err(),
				}
				return nil
			}

			defer func() {
				<-sem
			}()

			p, err := s.registry.Get(repo.Provider)
			if err != nil {
				results[i] = model.RepositoryResult{
					Repository: repo,
					Err:        err,
				}
				return nil
			}

			release, err := p.LatestRelease(ctx, repo)

			results[i] = model.RepositoryResult{
				Repository: repo,
				Release:    release,
				Err:        err,
			}

			return nil
		})
	}

	_ = g.Wait()

	return results
}
