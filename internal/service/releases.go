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

type Result struct {
	Index  int
	Result model.RepositoryResult
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
) <-chan Result {
	out := make(chan Result)

	go func() {
		defer close(out)

		g, ctx := errgroup.WithContext(ctx)
		sem := make(chan struct{}, s.concurrency)

		for i, repo := range repos {
			i, repo := i, repo

			g.Go(func() error {
				select {
				case sem <- struct{}{}:
				case <-ctx.Done():
					select {
					case out <- Result{
						Index: i,
						Result: model.RepositoryResult{
							Repository: repo,
							Err:        ctx.Err(),
						},
					}:
					case <-ctx.Done():
					}
					return nil
				}

				defer func() {
					<-sem
				}()

				p, err := s.registry.Get(repo.Provider)
				if err != nil {
					out <- Result{
						Index: i,
						Result: model.RepositoryResult{
							Repository: repo,
							Err:        err,
						},
					}
					return nil
				}

				release, err := p.LatestRelease(ctx, repo)

				out <- Result{
					Index: i,
					Result: model.RepositoryResult{
						Repository: repo,
						Release:    release,
						Err:        err,
					},
				}

				return nil
			})
		}

		_ = g.Wait()
	}()

	return out
}
