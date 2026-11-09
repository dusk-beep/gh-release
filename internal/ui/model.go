package ui

import (
	"context"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/table"

	"gh-release/internal/model"
	"gh-release/internal/service"
)

type Model struct {
	fetchCtx     context.Context
	service      *service.Service
	fetchCancel  context.CancelFunc
	fetchResults <-chan service.Result

	releaseMarkdown string
	releaseTitle    string
	repos           []model.Repository

	results []model.RepositoryResult

	displayOrder []int

	spinner spinner.Model

	table table.Model

	width  int
	height int

	releaseScroll int

	loading bool

	showRelease bool
}

func New(svc *service.Service, repos []model.Repository) Model {
	results := make([]model.RepositoryResult, len(repos))

	for i, repo := range repos {
		results[i] = model.RepositoryResult{
			Repository: repo,
		}
	}

	return Model{
		service: svc,
		repos:   repos,
		results: results,
		table:   newTable(),
		spinner: spinner.New(),
		loading: true,
	}
}
