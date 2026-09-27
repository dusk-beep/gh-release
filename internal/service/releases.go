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
				// Acquire a worker slot, but don't wait forever
				// if the fetch has been cancelled.
				select {
				case sem <- struct{}{}:
				case <-ctx.Done():
					return nil
				}

				defer func() {
					<-sem
				}()

				p, err := s.registry.Get(repo.Provider)
				if err != nil {
					return s.emit(ctx, out, Result{
						Index: i,
						Result: model.RepositoryResult{
							Repository: repo,
							Err:        err,
						},
					})
				}

				release, err := p.LatestRelease(ctx, repo)

				return s.emit(ctx, out, Result{
					Index: i,
					Result: model.RepositoryResult{
						Repository: repo,
						Release:    release,
						Err:        err,
					},
				})
			})
		}

		_ = g.Wait()
	}()

	return out
}

func (s *Service) emit(
	ctx context.Context,
	out chan<- Result,
	result Result,
) error {
	select {
	case out <- result:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
